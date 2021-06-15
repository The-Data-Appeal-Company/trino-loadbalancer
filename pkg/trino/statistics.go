package trino

import (
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/models"
)

type Retriever interface {
	ClusterStatistics(models.Coordinator) (models.ClusterStatistics, error)
	QueryDetail(coordinator models.Coordinator, queryID string) (models.QueryDetail, error)
	QueryList(coordinator models.Coordinator) (models.QueryList, error)
}
