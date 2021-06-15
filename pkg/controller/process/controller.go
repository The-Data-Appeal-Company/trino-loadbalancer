package process

import (
	"context"
	"fmt"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/api/trino"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/concurrency"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/models"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/discovery"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/proxy/healthcheck"
	"time"
)

type Controller struct {
	Api         trino.Api
	Discovery   discovery.Storage
	HealthCheck healthcheck.HealthCheck
	State       State
}

func NewController(api trino.Api, discovery discovery.Storage, healthCheck healthcheck.HealthCheck, state State) Controller {
	return Controller{
		Api:         api,
		Discovery:   discovery,
		HealthCheck: healthCheck,
		State:       state,
	}
}

func (c Controller) Run(ctx context.Context) error {
	coordinators, err := c.Discovery.All(ctx)
	if err != nil {
		return err
	}

	mg := concurrency.NewMultiErrorGroup()

	for _, coord := range coordinators {
		mg.Go(func() error {
			return c.controlCluster(ctx, coord)
		})
	}

	err = mg.Wait()
	if err != nil {
		return err
	}

	return nil
}

func (c Controller) controlCluster(ctx context.Context, cluster models.Coordinator) error {

	currentState := time.Now()
	previousState, err := c.State.Get(ctx, cluster)
	if err != nil {
		return err
	}

	health, err := c.HealthCheck.Check(cluster.URL)
	if err != nil {
		return err
	}

	if health.Status != healthcheck.StatusHealthy {
		return nil
	}

	queriesList, err := c.Api.QueryList(cluster)
	if err != nil {
		return err
	}

	completedQueryList := c.filterProcessedQueries(queriesList, previousState)

	for _, query := range completedQueryList {
		queryDetail, err := c.Api.QueryDetail(cluster, query.QueryId)
		if err != nil {
			return err
		}
		fmt.Println(queryDetail)
	}

	if err := c.State.Set(ctx, cluster, currentState); err != nil {
		return err
	}

	return nil
}

func (c Controller) filterProcessedQueries(list trino.QueryList, lastExecution time.Time) trino.QueryList {

	filterQueryList := make(trino.QueryList, 0)
	for _, item := range list {
		if item.State != trino.QueryFinished {
			continue
		}

		if lastExecution.After(item.Session.Start) {
			continue
		}

		filterQueryList = append(filterQueryList, item)
	}
	return filterQueryList
}
