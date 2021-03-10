package statistics

import (
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/models"
)

func Noop() NoOp {
	return NoOp{}
}

type NoOp struct {
}

func (n NoOp) GetStatistics(models.Coordinator) (models.ClusterStatistics, error) {
	return models.ClusterStatistics{}, nil
}
