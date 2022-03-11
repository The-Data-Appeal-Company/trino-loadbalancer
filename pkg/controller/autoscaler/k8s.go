package autoscaler

import (
	"context"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/api/trino"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/logging"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"net/url"
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
	Coordinator *url.URL
	Namespace   string
	Deployment  string
	Min         int
	Max         int
	ScaleAfter  time.Duration
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

	needScaleUp := hasQueriesInState(queries, StateWaitingForResources)
	if needScaleUp {
		k.logger.Info("found at least one query in waiting, trigger scale up to %d", request.Max)
		return k.scaleCluster(request.Namespace, request.Deployment, request.Max)
	}

	needScaleDown, err := k.needScaleDown(request, queries)
	if err != nil {
		return err
	}

	if needScaleDown {
		k.logger.Info("elapsed at least %s from last query, trigger scale down to %d", request.ScaleAfter, request.Min)
		return k.scaleCluster(request.Namespace, request.Deployment, request.Min)
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

func (k *KubeClientAutoscaler) scaleCluster(namespace string, deployment string, replicas int) error {
	ctx := context.TODO()

	scaleOpt := &autoscalingv1.Scale{
		ObjectMeta: v1.ObjectMeta{
			Name:      deployment,
			Namespace: namespace,
		},
		Spec: autoscalingv1.ScaleSpec{
			Replicas: int32(replicas),
		},
	}

	_, err := k.client.AppsV1().
		Deployments(namespace).
		UpdateScale(ctx, deployment, scaleOpt, v1.UpdateOptions{})

	if err != nil {
		return err
	}

	return nil
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
