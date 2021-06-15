package trino

import (
	models2 "github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/models"
)

func Mock(statistics models2.ClusterStatistics, err error) MockStats {
	return MockStats{
		statistics: statistics,
		err:        err,
	}
}

type MockStats struct {
	statistics models2.ClusterStatistics
	err        error
}

func (m MockStats) ClusterStatistics(models2.Coordinator) (models2.ClusterStatistics, error) {
	return m.statistics, m.err
}

func (m MockStats) QueryDetail(coord models2.Coordinator, queryID string) (models2.QueryDetail, error) {
	return models2.QueryDetail{}, nil
}

func (m MockStats) QueryList(coord models2.Coordinator) (models2.QueryList, error) {
	return nil, nil
}
