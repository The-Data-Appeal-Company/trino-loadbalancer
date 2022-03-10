package trino

import (
	"net/url"
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

func (m MockStats) ClusterStatistics(*url.URL) (ClusterStatistics, error) {
	return m.statistics, m.err
}

func (m MockStats) QueryDetail(coord *url.URL, queryID string) (QueryDetail, error) {
	return QueryDetail{}, nil
}

func (m MockStats) QueryList(coord *url.URL) (QueryList, error) {
	return nil, nil
}

type MockApi struct {
	ClusterStatisticsFn func(*url.URL) (ClusterStatistics, error)
	QueryDetailFn       func(coordinator *url.URL, queryID string) (QueryDetail, error)
	QueryListFn         func(coordinator *url.URL) (QueryList, error)
}

func (m MockApi) ClusterStatistics(coordinator *url.URL) (ClusterStatistics, error) {
	if m.ClusterStatisticsFn == nil {
		return ClusterStatistics{}, nil
	}
	return m.ClusterStatisticsFn(coordinator)
}

func (m MockApi) QueryDetail(coordinator *url.URL, queryID string) (QueryDetail, error) {
	if m.QueryDetailFn == nil {
		return QueryDetail{}, nil
	}
	return m.QueryDetailFn(coordinator, queryID)
}

func (m MockApi) QueryList(coordinator *url.URL) (QueryList, error) {
	if m.QueryListFn == nil {
		return QueryList{}, nil
	}
	return m.QueryListFn(coordinator)
}
