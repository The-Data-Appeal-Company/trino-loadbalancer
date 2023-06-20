package autoscaler

import (
	"context"
	"errors"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/api/trino"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/logging"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/configuration"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"net/url"
	"regexp"
	"time"
)

const (
	StateWaitingForResources string = "WAITING_FOR_RESOURCES"
	StateRunning             string = "RUNNING"
)

type KubeAutoscaler interface {
	Execute(KubeRequest) error
}

type KubeRequest struct {
	Coordinator  *url.URL
	Namespace    string
	Deployment   string
	Min          int
	Max          int
	ScaleAfter   time.Duration
	DynamicScale configuration.AutoscalerDynamicScale
}

type KubeClientAutoscaler struct {
	client   kubernetes.Interface
	trinoApi trino.Api
	state    State
	logger   logging.Logger
}

func NewKubeClientAutoscaler(client kubernetes.Interface, trinoApi trino.Api, state State, logger logging.Logger) *KubeClientAutoscaler {
	return &KubeClientAutoscaler{client: client, trinoApi: trinoApi, state: state, logger: logger}
}

func (k *KubeClientAutoscaler) Execute(request KubeRequest) error {
	queries, err := k.trinoApi.QueryList(request.Coordinator)
	if err != nil {
		return err
	}

	if !request.DynamicScale.Enabled {
		needScaleUp := hasQueriesInState(queries, StateWaitingForResources)
		if needScaleUp {
			k.logger.Info("found at least one query in waiting, trigger scale up to %d", request.Max)
			return k.scaleCluster(request, request.Max)
		}
	} else {
		needScaleUp, instances, err := k.needScaleUp(request, queries)
		if err != nil {
			return err
		}
		if needScaleUp {
			return k.scaleCluster(request, instances)
		}
	}

	needScaleDown, err := k.needScaleDown(request, queries)
	if err != nil {
		return err
	}

	if needScaleDown {
		k.logger.Info("elapsed at least %s from last query, trigger scale down to %d", request.ScaleAfter, request.Min)
		return k.scaleCluster(request, request.Min)
	}

	return nil
}

func (k *KubeClientAutoscaler) needScaleDown(req KubeRequest, queries trino.QueryList) (bool, error) {
	hasRunningQueries := hasQueriesInState(queries, StateRunning)
	if hasRunningQueries {
		return false, nil
	}

	var lastQueryTime time.Time

	// if the trino api doesn't return any query we fetch the last execution info from the state
	if len(queries) == 0 {
		lastState, err := k.state.LastQueryExecution(req.Coordinator.String())
		if err != nil {
			return false, err
		}
		lastQueryTime = lastState
	} else {
		// if the api has returned atleast one query we use the most recent query end time
		lastQueryTime = lastQueryExecution(queries)
	}

	if lastQueryTime.IsZero() {
		lastQueryTime = time.Now()
	}

	// we save the result into the state, if no queries are found and no previous state was set we set Now() as last time
	if err := k.state.SetLastQueryExecution(req.Coordinator.String(), lastQueryTime); err != nil {
		return false, err
	}

	k.logger.Info("time pass since last query %s, need to scale down: %t", time.Since(lastQueryTime), time.Since(lastQueryTime) > req.ScaleAfter)
	return time.Since(lastQueryTime) > req.ScaleAfter, nil
}

func (k *KubeClientAutoscaler) needScaleUp(request KubeRequest, queries trino.QueryList) (bool, int, error) {
	runningQuery := filterByState(queries, StateRunning)
	waitingQuery := filterByState(queries, StateWaitingForResources)

	allQueries := append(runningQuery, waitingQuery...)
	if len(allQueries) == 0 {
		return false, 0, nil
	}

	lastInstances, err := k.state.GetClusterInstances(request.Coordinator.String())
	if err != nil {
		return false, 0, err
	}
	if err != nil {
		if errors.Is(err, NoInstancesInStateError) {
			lastInstances, err = k.getDeploymentInstances(request.Namespace, request.Deployment)
			if err != nil {
				return false, 0, err
			}
		} else {
			return false, 0, err
		}
	}
	scaleInstances := request.DynamicScale.Default

	for _, rule := range request.DynamicScale.Rules {
		r, err := regexp.Compile(rule.Regexp)
		if err != nil {
			k.logger.Warn("cannot parse regexp '%s' of dynamic rule", rule.Regexp)
			break
		}
		for _, query := range queries {
			if r.MatchString(query.Session.User) {
				scaleInstances = maxInt(scaleInstances, rule.Instances)
				break
			}
		}
	}

	return lastInstances < int32(scaleInstances),
		scaleInstances,
		nil
}

func (k *KubeClientAutoscaler) scaleCluster(request KubeRequest, replicas int) error {
	ctx := context.TODO()

	scaleOpt := &autoscalingv1.Scale{
		ObjectMeta: v1.ObjectMeta{
			Name:      request.Deployment,
			Namespace: request.Namespace,
		},
		Spec: autoscalingv1.ScaleSpec{
			Replicas: int32(replicas),
		},
	}

	_, err := k.client.AppsV1().
		Deployments(request.Namespace).
		UpdateScale(ctx, request.Deployment, scaleOpt, v1.UpdateOptions{})

	if err != nil {
		return err
	}

	err = k.state.SetClusterInstances(request.Coordinator.String(), int32(replicas))
	if err != nil {
		return err
	}

	return nil
}

func (k *KubeClientAutoscaler) getDeploymentInstances(namespace string, deployment string) (int32, error) {
	ctx := context.TODO()

	info, err := k.client.AppsV1().
		Deployments(namespace).
		Get(ctx, deployment, v1.GetOptions{})

	if err != nil {
		return 0, err
	}

	if info.Spec.Replicas != nil {
		return *info.Spec.Replicas, nil
	}

	return 0, nil
}

func filterByState(queries trino.QueryList, state string) trino.QueryList {
	results := make(trino.QueryList, 0)
	for _, query := range queries {
		if query.State == state {
			results = append(results, query)
		}
	}
	return results
}

func hasQueriesInState(queries trino.QueryList, state string) bool {
	for _, query := range queries {
		if query.State == state {
			return true
		}
	}
	return false
}

func lastQueryExecution(queries trino.QueryList) time.Time {
	var last = time.Time{}
	for _, query := range queries {
		endTime := query.QueryStats.EndTime
		if endTime.IsZero() {
			continue
		}

		if last.IsZero() || endTime.After(last) {
			last = endTime
		}
	}
	return last
}

func maxInt(value1, value2 int) int {
	if value1 > value2 {
		return value1
	}
	return value2
}
