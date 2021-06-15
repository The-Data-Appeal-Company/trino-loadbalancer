package trino

import (
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/models"
)

func Mock(statistics ClusterStatistics, err error) MockStats {
	return MockStats{
		statistics: statistics,
		err:        err,
	}
}

type MockStats struct {
	statistics ClusterStatistics
	err        error
}

func (m MockStats) ClusterStatistics(models.Coordinator) (ClusterStatistics, error) {
	return m.statistics, m.err
}

func (m MockStats) QueryDetail(coord models.Coordinator, queryID string) (QueryDetail, error) {
	return QueryDetail{}, nil
}

func (m MockStats) QueryList(coord models.Coordinator) (QueryList, error) {
	return nil, nil
}
