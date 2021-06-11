package statistics

import (
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/models"
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

func (m MockStats) ClusterStatistics(models.Coordinator) (models.ClusterStatistics, error) {
	return m.statistics, m.err
}

func (m MockStats) QueryStatistics(coord models.Coordinator, queryID string) (models.QueryStats, error) {
	return models.QueryStats{}, nil
}
