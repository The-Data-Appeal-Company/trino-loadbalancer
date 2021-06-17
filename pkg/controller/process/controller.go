package process

import (
	"context"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/api/trino"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/concurrency"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/models"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/controller/components"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/discovery"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/proxy/healthcheck"
	"time"
)

type Controller struct {
	api          trino.Api
	discovery    discovery.Storage
	healthCheck  healthcheck.HealthCheck
	state        State
	queryHandler components.QueryHandler
}

func NewController(api trino.Api, discovery discovery.Storage, healthCheck healthcheck.HealthCheck, state State, queryHandler components.QueryHandler) Controller {
	return Controller{api: api, discovery: discovery, healthCheck: healthCheck, state: state, queryHandler: queryHandler}
}

func (c Controller) Run(ctx context.Context) error {
	coordinators, err := c.discovery.All(ctx)
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
	previousState, err := c.state.Get(ctx, cluster)
	if err != nil {
		return err
	}

	health, err := c.healthCheck.Check(cluster.URL)
	if err != nil {
		return err
	}

	if health.Status != healthcheck.StatusHealthy {
		return nil
	}

	queriesList, err := c.api.QueryList(cluster)
	if err != nil {
		return err
	}

	completedQueryList := c.filterProcessedQueries(queriesList, previousState)

	for _, query := range completedQueryList {
		queryDetail, err := c.api.QueryDetail(cluster, query.QueryId)
		if err != nil {
			return err
		}

		if err := c.queryHandler.Execute(ctx, queryDetail); err != nil {
			return err
		}
	}

	if err := c.state.Set(ctx, cluster, currentState); err != nil {
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

		if lastExecution.After(item.QueryStats.CreateTime) {
			continue
		}

		filterQueryList = append(filterQueryList, item)
	}
	return filterQueryList
}
