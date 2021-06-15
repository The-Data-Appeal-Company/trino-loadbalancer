package trino

import (
	models2 "github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/models"
)

type Retriever interface {
	ClusterStatistics(models2.Coordinator) (models2.ClusterStatistics, error)
	QueryDetail(coordinator models2.Coordinator, queryID string) (models2.QueryDetail, error)
	QueryList(coordinator models2.Coordinator) (models2.QueryList, error)
}
