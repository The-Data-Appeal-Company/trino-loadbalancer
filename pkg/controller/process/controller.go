package process

import (
	"context"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/api/trino"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/concurrency"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/logging"
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
	logger       logging.Logger
}

func NewController(api trino.Api, discovery discovery.Storage, healthCheck healthcheck.HealthCheck, state State, queryHandler components.QueryHandler, logger logging.Logger) Controller {
	return Controller{api: api, discovery: discovery, healthCheck: healthCheck, state: state, queryHandler: queryHandler, logger: logger}
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
	previousState, err := c.state.Get(ctx, cluster)
	if err != nil {
		return err
	}

	health, err := c.healthCheck.Check(cluster.URL)
	if err != nil {
		return err
	}

	if health.Status != healthcheck.StatusHealthy {
		c.logger.Info("skipping unhealthy cluster: %s ( %s )", cluster.Name, cluster.URL.String())
		return nil
	}

	queriesList, err := c.api.QueryList(cluster)
	if err != nil {
		return err
	}

	completedQueryList := c.filterProcessedQueries(queriesList, previousState)

	if len(completedQueryList) == 0 {
		return nil
	}

	c.logger.Info("%s: retrieved %d queries", cluster.Name, len(completedQueryList))

	for _, query := range completedQueryList {
		queryDetail, err := c.api.QueryDetail(cluster, query.QueryId)
		if err != nil {
			if err == trino.ErrQueryNotFound {
				continue
			}
			return err
		}

		if err := c.queryHandler.Execute(ctx, queryDetail); err != nil {
			return err
		}
	}

	if lastQueryTime := lastQueryTime(completedQueryList); !lastQueryTime.IsZero() {
		if err := c.state.Set(ctx, cluster, lastQueryTime); err != nil {
			return err
		}
	}

	return nil
}

func lastQueryTime(list trino.QueryList) time.Time {
	var max time.Time
	for _, q := range list {
		if q.QueryStats.CreateTime.After(max) {
			max = q.QueryStats.CreateTime
		}
	}
	return max
}

func (c Controller) filterProcessedQueries(list trino.QueryList, lastExecution time.Time) trino.QueryList {
	filterQueryList := make(trino.QueryList, 0)

	for _, item := range list {
		if item.State != trino.QueryFinished {
			continue
		}

		if item.QueryStats.TotalDrivers < 2 {
			continue
		}

		if item.QueryStats.CreateTime.Before(lastExecution) || item.QueryStats.CreateTime.Equal(lastExecution) {
			continue
		}

		filterQueryList = append(filterQueryList, item)
	}

	return filterQueryList
}
