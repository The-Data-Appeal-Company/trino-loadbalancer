package factory

import (
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/healthcheck"
)

type HealthCheckConfiguration struct {
	Enabled bool
}

func CreateHealthCheck(conf HealthCheckConfiguration) (healthcheck.HealthCheck, error) {
	if !conf.Enabled {
		return healthcheck.NoOp(), nil
	}

	return healthcheck.NewHttpHealth(), nil
}
