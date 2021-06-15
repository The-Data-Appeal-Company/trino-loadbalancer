package trino

import (
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/models"
)

func Noop() NoOp {
	return NoOp{}
}

type NoOp struct {
}

func (n NoOp) ClusterStatistics(models.Coordinator) (models.ClusterStatistics, error) {
	return models.ClusterStatistics{}, nil
}

func (n NoOp) QueryDetail(coord models.Coordinator, queryID string) (models.QueryDetail, error) {
	return models.QueryDetail{}, nil
}
func (n NoOp) QueryList(coord models.Coordinator) (models.QueryList, error) {
	return nil, nil
}
