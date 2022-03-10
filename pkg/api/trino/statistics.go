package trino

import (
	"net/url"
)

type Api interface {
	ClusterStatistics(url *url.URL) (ClusterStatistics, error)
	QueryDetail(url *url.URL, queryID string) (QueryDetail, error)
	QueryList(url *url.URL) (QueryList, error)
}
