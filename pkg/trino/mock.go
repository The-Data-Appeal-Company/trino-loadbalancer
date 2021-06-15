package trino

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

func (m MockStats) QueryDetail(coord models.Coordinator, queryID string) (models.QueryDetail, error) {
	return models.QueryDetail{}, nil
}

func (m MockStats) QueryList(coord models.Coordinator) (models.QueryList, error) {
	return nil, nil
}
