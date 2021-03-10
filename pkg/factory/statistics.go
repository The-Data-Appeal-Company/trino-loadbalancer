package factory

import (
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/statistics"
)

type StatisticsConfiguration struct {
	Enabled bool
}

func CreateStatisticsRetriever(conf StatisticsConfiguration) (statistics.Retriever, error) {
	if !conf.Enabled {
		return statistics.Noop(), nil
	}

	return statistics.NewPrestoClusterApi(), nil
}
