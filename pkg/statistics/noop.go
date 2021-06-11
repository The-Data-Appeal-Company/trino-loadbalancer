package statistics

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

func (n NoOp) QueryStatistics(coord models.Coordinator, queryID string) (models.QueryStats, error) {
	return models.QueryStats{}, nil
}
