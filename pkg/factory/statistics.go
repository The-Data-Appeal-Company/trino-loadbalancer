package factory

import (
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/trino"
)

type StatisticsConfiguration struct {
	Enabled bool
}

func CreateStatisticsRetriever(conf StatisticsConfiguration) (trino.Retriever, error) {
	if !conf.Enabled {
		return trino.Noop(), nil
	}

	return trino.NewClusterApi(), nil
}
