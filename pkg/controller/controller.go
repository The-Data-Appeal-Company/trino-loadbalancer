package controller

import (
	"context"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/api/trino"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/concurrency"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/models"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/discovery"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/proxy/healthcheck"
)

type Controller struct {
	Api         trino.Api
	Discovery   discovery.Storage
	HealthCheck healthcheck.HealthCheck
}

func NewController(api trino.Api, discovery discovery.Storage, healthCheck healthcheck.HealthCheck) Controller {
	return Controller{
		Api:       api,
		Discovery: discovery,
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
			return c.controlCluster(coord)
		})
	}

	err = mg.Wait()
	if err != nil {
		return err
	}

	return nil
}

func (c Controller) controlCluster(cluster models.Coordinator) error {
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

	queriesList = c.filterProcessedQueries(queriesList)

	for _, query := range queriesList {
		queryDetail, err := c.Api.QueryDetail(cluster, query.QueryId)
		if err != nil {
			return err
		}


	}


	return nil
}

func (c Controller) filterProcessedQueries(list trino.QueryList) trino.QueryList {
	return list
}
