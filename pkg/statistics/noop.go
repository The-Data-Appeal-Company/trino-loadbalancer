package statistics

import (
	"github.com/The-Data-Appeal-Company/presto-loadbalancer/pkg/models"
)

func Noop() NoOp {
	return NoOp{}
}

type NoOp struct {
}

func (n NoOp) GetStatistics(models.Coordinator) (models.ClusterStatistics, error) {
	return models.ClusterStatistics{}, nil
}
