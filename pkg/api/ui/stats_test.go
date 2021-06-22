package ui

import (
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/api/trino"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/logging"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/discovery"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStatsApi(t *testing.T) {
	stats := trino.Mock(trino.ClusterStatistics{}, nil)
	discover := discovery.Noop()
	discoverStorage := discovery.NewMemoryStorage()

	api := NewApi(stats, discover, discoverStorage, logging.Noop())

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/stats", nil)

	api.statistics(rr, req)

	require.Equal(t, rr.Code, http.StatusOK)
}
