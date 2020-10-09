package statistics

import (
	"github.com/The-Data-Appeal-Company/presto-loadbalancer/pkg/models"
)

func Mock(statistics models.ClusterStatistics, err error) MockStats {
	return MockStats{
		statistics: statistics,
		err:        err,
	}
}

type MockStats struct {
	statistics models.ClusterStatistics
	err        error
}

func (m MockStats) GetStatistics(models.Coordinator) (models.ClusterStatistics, error) {
	return m.statistics, m.err
}
