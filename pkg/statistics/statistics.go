package statistics

import (
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/models"
)

type Retriever interface {
	GetStatistics(models.Coordinator) (models.ClusterStatistics, error)
}
