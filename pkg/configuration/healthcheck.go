package configuration

import (
	"fmt"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/proxy/healthcheck"
)

type HealthCheckConfiguration struct {
	Enabled bool
	Type    string
}

const (
	healthCheckTypeQuery string = "query"
	healthCheckTypeHttp  string = "http"
)

func CreateHealthCheck(conf HealthCheckConfiguration) (healthcheck.HealthCheck, error) {
	if !conf.Enabled {
		return healthcheck.NoOp(), nil
	}

	return getHealthCheckFromType(conf.Type)
}

func getHealthCheckFromType(healthType string) (healthcheck.HealthCheck, error) {

	switch healthType {
	case healthCheckTypeQuery:
		return healthcheck.NewTrinoQueryHealth(), nil
	case healthCheckTypeHttp:
		return healthcheck.NewHttpHealth(), nil
	default:
		return healthcheck.NoOp(), fmt.Errorf("invalid health check type")
	}
}
