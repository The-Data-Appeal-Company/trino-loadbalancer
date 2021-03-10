package lb

import (
	"errors"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/healthcheck"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/logging"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/models"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/session"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/statistics"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func PoolConfigTest() PoolConfig {
	return PoolConfig{
		HealthCheckDelay: 5 * time.Second,
		StatisticsDelay:  5 * time.Second,
	}
}

func TestPool_AddHealthyBackend(t *testing.T) {

	sessStore := session.NewMemoryStorage()
	hc := healthcheck.Mock(healthcheck.Health{
		Status:    healthcheck.StatusHealthy,
		Message:   "ok",
		Timestamp: time.Now(),
	}, nil)

	stats := statistics.Mock(models.ClusterStatistics{}, nil)
	logger := logging.Noop()

	pool := NewPool(PoolConfigTest(), sessStore, hc, stats, logger)

	coord := models.Coordinator{
		Name:    "coord-0",
		URL:     mustUrl("http://presto.local:8080"),
		Tags:    map[string]string{},
		Enabled: true,
	}

	err := pool.Add(coord)
	require.NoError(t, err)

	// force pool health update before further checking
	// at the moment the health check is executed synchronously in the Add() method but this behaviour may change
	// for the purpose of this test we don't care about the Add method behaviour
	err = pool.UpdateStatus()
	require.NoError(t, err)

	backends := pool.AllBackends()
	require.Len(t, backends, 1)

	availables, err := pool.AvailableBackends()
	require.NoError(t, err)
	require.Len(t, availables, len(backends))

	first := availables[0]
	require.Equal(t, first.Backend, coord)

}

func TestPool_AddUnHealthyBackend(t *testing.T) {

	sessStore := session.NewMemoryStorage()
	hc := healthcheck.Mock(healthcheck.Health{
		Status:    healthcheck.StatusUnhealthy,
		Message:   "generic health check error",
		Timestamp: time.Now(),
	}, nil)

	stats := statistics.Mock(models.ClusterStatistics{}, nil)
	logger := logging.Noop()

	pool := NewPool(PoolConfigTest(), sessStore, hc, stats, logger)

	coord := models.Coordinator{
		Name:    "coord-0",
		URL:     mustUrl("http://presto.local:8080"),
		Tags:    map[string]string{},
		Enabled: true,
	}

	err := pool.Add(coord)
	require.NoError(t, err)

	err = pool.UpdateStatus()
	require.NoError(t, err)

	backends := pool.AllBackends()
	require.Len(t, backends, 1)

	_, err = pool.AvailableBackends()
	require.Error(t, ErrNoBackendsAvailable)

}

func TestPool_RemoveBackend(t *testing.T) {

	sessStore := session.NewMemoryStorage()
	hc := healthcheck.Mock(healthcheck.Health{
		Status:    healthcheck.StatusHealthy,
		Message:   "ok",
		Timestamp: time.Now(),
	}, nil)

	stats := statistics.Mock(models.ClusterStatistics{}, nil)
	logger := logging.Noop()

	pool := NewPool(PoolConfigTest(), sessStore, hc, stats, logger)

	err := pool.Add(models.Coordinator{
		Name:    "coord-0",
		URL:     mustUrl("http://presto.local:8080"),
		Tags:    map[string]string{},
		Enabled: true,
	})
	require.NoError(t, err)

	err = pool.Add(models.Coordinator{
		Name:    "coord-1",
		URL:     mustUrl("http://presto.local:8080"),
		Tags:    map[string]string{},
		Enabled: true,
	})
	require.NoError(t, err)

	// force pool health update before further checking
	// at the moment the health check is executed synchronously in the Add() method but this behaviour may change
	// for the purpose of this test we don't care about the Add method behaviour
	err = pool.UpdateStatus()
	require.NoError(t, err)

	backends := pool.AllBackends()
	require.Len(t, backends, 2)

	err = pool.Remove("coord-0")
	require.NoError(t, err)

	backends = pool.AllBackends()
	require.Len(t, backends, 1)

	_, err = pool.GetByName("coord-1", healthcheck.StatusHealthy)
	require.NoError(t, err)
}

func TestPool_GetByName(t *testing.T) {

	sessStore := session.NewMemoryStorage()
	hc := healthcheck.Mock(healthcheck.Health{
		Status:    healthcheck.StatusHealthy,
		Message:   "ok",
		Timestamp: time.Now(),
	}, nil)

	stats := statistics.Mock(models.ClusterStatistics{}, nil)
	logger := logging.Noop()

	pool := NewPool(PoolConfigTest(), sessStore, hc, stats, logger)

	_, err := pool.GetByName("test", healthcheck.StatusHealthy)
	require.True(t, errors.Is(err, ErrNoBackendsAvailable))
}

func TestPool_GetByNameWithUnhealthyStatus(t *testing.T) {
	sessStore := session.NewMemoryStorage()
	hc := healthcheck.Mock(healthcheck.Health{
		Status:    healthcheck.StatusUnhealthy,
		Message:   "ok",
		Timestamp: time.Now(),
	}, nil)

	stats := statistics.Mock(models.ClusterStatistics{}, nil)
	logger := logging.Noop()

	pool := NewPool(PoolConfigTest(), sessStore, hc, stats, logger)
	err := pool.Add(models.Coordinator{
		Name:    "coord-1",
		URL:     mustUrl("http://presto.local:8080"),
		Tags:    map[string]string{},
		Enabled: true,
	})
	require.NoError(t, err)

	err = pool.UpdateStatus()
	require.NoError(t, err)

	_, err = pool.GetByName("coord-1", healthcheck.StatusUnhealthy)
	require.NoError(t, err)
}

func TestPool_UpdateBackend(t *testing.T) {

	sessStore := session.NewMemoryStorage()
	hc := healthcheck.Mock(healthcheck.Health{
		Status:    healthcheck.StatusHealthy,
		Message:   "ok",
		Timestamp: time.Now(),
	}, nil)

	stats := statistics.Mock(models.ClusterStatistics{}, nil)
	logger := logging.Noop()

	pool := NewPool(PoolConfigTest(), sessStore, hc, stats, logger)

	coord := models.Coordinator{
		Name:    "coord-0",
		URL:     mustUrl("http://presto.local:8080"),
		Tags:    map[string]string{},
		Enabled: true,
	}

	err := pool.Add(coord)
	require.NoError(t, err)

	state, err := pool.GetByName(coord.Name, healthcheck.StatusHealthy)
	require.NoError(t, err)

	require.Equal(t, state.Backend, coord)

	newState := models.Coordinator{
		Tags: map[string]string{
			"updated": "true",
		},
		Enabled: false,
	}
	err = pool.Update(coord.Name, newState)
	require.NoError(t, err)

	require.Equal(t, state.Backend, models.Coordinator{
		Name:    coord.Name,
		URL:     coord.URL,
		Tags:    newState.Tags,
		Enabled: newState.Enabled,
	})

}
