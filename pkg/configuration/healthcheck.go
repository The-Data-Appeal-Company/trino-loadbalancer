package configuration

import (
	healthcheck2 "github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/proxy/healthcheck"
)

type HealthCheckConfiguration struct {
	Enabled bool
}

func CreateHealthCheck(conf HealthCheckConfiguration) (healthcheck2.HealthCheck, error) {
	if !conf.Enabled {
		return healthcheck2.NoOp(), nil
	}

	return healthcheck2.NewHttpHealth(), nil
}
