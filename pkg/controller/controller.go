package controller

import (
	"context"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/api/trino"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/models"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/discovery"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/proxy/healthcheck"
	"golang.org/x/sync/errgroup"
	"sync"
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

	g, ctx := errgroup.WithContext(ctx)


	return nil
}

func controlCluster(cluster models.Coordinator) error {

	return nil
}
