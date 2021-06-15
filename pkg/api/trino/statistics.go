package trino

import (
	models2 "github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/models"
)

type Api interface {
	ClusterStatistics(models2.Coordinator) (ClusterStatistics, error)
	QueryDetail(coordinator models2.Coordinator, queryID string) (QueryDetail, error)
	QueryList(coordinator models2.Coordinator) (QueryList, error)
}
