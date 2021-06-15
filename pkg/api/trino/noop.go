package trino

import (
	models2 "github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/models"
)

func Noop() NoOp {
	return NoOp{}
}

type NoOp struct {
}

func (n NoOp) ClusterStatistics(models2.Coordinator) (ClusterStatistics, error) {
	return ClusterStatistics{}, nil
}

func (n NoOp) QueryDetail(coord models2.Coordinator, queryID string) (QueryDetail, error) {
	return QueryDetail{}, nil
}
func (n NoOp) QueryList(coord models2.Coordinator) (QueryList, error) {
	return nil, nil
}
