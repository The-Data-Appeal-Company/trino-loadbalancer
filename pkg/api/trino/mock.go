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

type MockApi struct {
	ClusterStatisticsFn func(models.Coordinator) (ClusterStatistics, error)
	QueryDetailFn       func(coordinator models.Coordinator, queryID string) (QueryDetail, error)
	QueryListFn         func(coordinator models.Coordinator) (QueryList, error)
}

func (m MockApi) ClusterStatistics(coordinator models.Coordinator) (ClusterStatistics, error) {
	if m.ClusterStatisticsFn == nil {
		return ClusterStatistics{}, nil
	}
	return m.ClusterStatisticsFn(coordinator)
}

func (m MockApi) QueryDetail(coordinator models.Coordinator, queryID string) (QueryDetail, error) {
	if m.QueryDetailFn == nil {
		return QueryDetail{}, nil
	}
	return m.QueryDetailFn(coordinator, queryID)
}

func (m MockApi) QueryList(coordinator models.Coordinator) (QueryList, error) {
	if m.QueryListFn == nil {
		return QueryList{}, nil
	}
	return m.QueryListFn(coordinator)
}
