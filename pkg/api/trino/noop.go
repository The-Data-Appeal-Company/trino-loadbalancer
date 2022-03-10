package trino

import (
	"net/url"
)

func Noop() NoOp {
	return NoOp{}
}

type NoOp struct {
}

func (n NoOp) ClusterStatistics(*url.URL) (ClusterStatistics, error) {
	return ClusterStatistics{}, nil
}

func (n NoOp) QueryDetail(coord *url.URL, queryID string) (QueryDetail, error) {
	return QueryDetail{}, nil
}
func (n NoOp) QueryList(coord *url.URL) (QueryList, error) {
	return nil, nil
}
