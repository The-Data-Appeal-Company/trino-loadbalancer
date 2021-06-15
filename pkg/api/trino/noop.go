package trino

import (
	models2 "github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/models"
)

func Noop() NoOp {
	return NoOp{}
}

type NoOp struct {
}

func (n NoOp) ClusterStatistics(models2.Coordinator) (models2.ClusterStatistics, error) {
	return models2.ClusterStatistics{}, nil
}

func (n NoOp) QueryDetail(coord models2.Coordinator, queryID string) (models2.QueryDetail, error) {
	return models2.QueryDetail{}, nil
}
func (n NoOp) QueryList(coord models2.Coordinator) (models2.QueryList, error) {
	return nil, nil
}
