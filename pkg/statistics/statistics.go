package statistics

import (
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/models"
)

type Retriever interface {
	ClusterStatistics(models.Coordinator) (models.ClusterStatistics, error)
	QueryStatistics(coord models.Coordinator, queryID string) (models.QueryStats, error)
}
