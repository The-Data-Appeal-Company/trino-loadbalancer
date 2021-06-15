package trino

import (
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/models"
)

type Api interface {
	ClusterStatistics(models.Coordinator) (ClusterStatistics, error)
	QueryDetail(coordinator models.Coordinator, queryID string) (QueryDetail, error)
	QueryList(coordinator models.Coordinator) (QueryList, error)
}
