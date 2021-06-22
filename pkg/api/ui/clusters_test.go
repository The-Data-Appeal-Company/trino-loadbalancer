package ui

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/api/trino"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/logging"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/models"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/tests"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/discovery"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClustersListApi(t *testing.T) {
	stats := trino.Mock(trino.ClusterStatistics{}, nil)
	discover := discovery.Noop()
	discoverStorage := discovery.NewMemoryStorage()

	err := discoverStorage.Add(context.TODO(), models.Coordinator{
		Name: "cluster-00",
		URL:  tests.MustUrl("http://localhost:8080"),
		Tags: map[string]string{
			"test": "true",
		},
		Enabled: true,
	})
	require.NoError(t, err)

	err = discoverStorage.Add(context.TODO(), models.Coordinator{
		Name: "cluster-01",
		URL:  tests.MustUrl("http://localhost:8081"),
		Tags: map[string]string{
			"test": "true",
		},
		Enabled: false,
	})
	require.NoError(t, err)

	api := NewApi(stats, discover, discoverStorage, logging.Noop())

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/stats", nil)

	api.clustersList(rr, req)

	require.Equal(t, rr.Code, http.StatusOK)

	var response []Cluster
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}

	require.Len(t, response, 2)
}

func TestAddClusterApi(t *testing.T) {
	stats := trino.Mock(trino.ClusterStatistics{}, nil)
	discover := discovery.Noop()
	discoverStorage := discovery.NewMemoryStorage()

	err := discoverStorage.Add(context.TODO(), models.Coordinator{
		Name: "cluster-00",
		URL:  tests.MustUrl("http://localhost:8080"),
		Tags: map[string]string{
			"test": "true",
		},
		Enabled: true,
	})
	require.NoError(t, err)

	api := NewApi(stats, discover, discoverStorage, logging.Noop())

	rr := httptest.NewRecorder()

	addReq := ClusterAddRequest{
		Name:    "cluster-01",
		Url:     "http://localhost:8082",
		Enabled: true,
	}

	addReqBody, err := json.Marshal(addReq)
	require.NoError(t, err)

	api.addCluster(rr, httptest.NewRequest(http.MethodGet, "http://localhost:8080/stats", bytes.NewBuffer(addReqBody)))

	require.Equal(t, rr.Code, http.StatusOK)

	clusters, err := discoverStorage.All(context.TODO())
	require.NoError(t, err)

	require.Len(t, clusters, 2)
}

func TestLaunchDiscoverApi(t *testing.T) {
	stats := trino.Mock(trino.ClusterStatistics{}, nil)
	discover := discovery.Noop()
	discoverStorage := discovery.NewMemoryStorage()

	api := NewApi(stats, discover, discoverStorage, logging.Noop())

	rr := httptest.NewRecorder()
	api.launchDiscover(rr, httptest.NewRequest(http.MethodGet, "http://localhost:8080/stats", nil))

	require.Equal(t, rr.Code, http.StatusOK)
}
