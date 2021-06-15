package configuration

import (
	trino2 "github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/api/trino"
)

type StatisticsConfiguration struct {
	Enabled bool
}

func CreateStatisticsRetriever(conf StatisticsConfiguration) (trino2.Api, error) {
	if !conf.Enabled {
		return trino2.Noop(), nil
	}

	return trino2.NewClusterApi(), nil
}
