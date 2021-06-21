package configuration

import (
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/api/trino"
)

type StatisticsConfiguration struct {
	Enabled bool
}

func CreateStatisticsRetriever(conf StatisticsConfiguration) (trino.Api, error) {
	if !conf.Enabled {
		return trino.Noop(), nil
	}

	return trino.NewClusterApi(), nil
}
