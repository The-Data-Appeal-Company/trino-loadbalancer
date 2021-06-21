package trino

import (
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/models"
)

func Noop() NoOp {
	return NoOp{}
}

type NoOp struct {
}

func (n NoOp) ClusterStatistics(models.Coordinator) (ClusterStatistics, error) {
	return ClusterStatistics{}, nil
}

func (n NoOp) QueryDetail(coord models.Coordinator, queryID string) (QueryDetail, error) {
	return QueryDetail{}, nil
}
func (n NoOp) QueryList(coord models.Coordinator) (QueryList, error) {
	return nil, nil
}
