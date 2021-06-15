package lb

import (
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/healthcheck"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/logging"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/models"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/session"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/trino"
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

	stats := trino.Mock(models.ClusterStatistics{}, nil)
	logger := logging.Noop()

	pool := NewPool(PoolConfigTest(), sessStore, hc, stats, logger)

	coord := models.Coordinator{
		Name:    "coord-0",
		URL:     mustUrl("http://trino.local:8080"),
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

	backends := pool.Fetch(FetchRequest{})
	require.Len(t, backends, 1)

	availables := pool.Fetch(FetchRequest{
		Status: ClusterStatusEnabled,
	})
	require.NoError(t, err)
	require.Len(t, availables, len(backends))

	first := availables[0]
	require.Equal(t, first.Coordinator, coord)

}

func TestPool_AddUnHealthyBackend(t *testing.T) {

	sessStore := session.NewMemoryStorage()
	hc := healthcheck.Mock(healthcheck.Health{
		Status:    healthcheck.StatusUnhealthy,
		Message:   "generic health check error",
		Timestamp: time.Now(),
	}, nil)

	stats := trino.Mock(models.ClusterStatistics{}, nil)
	logger := logging.Noop()

	pool := NewPool(PoolConfigTest(), sessStore, hc, stats, logger)

	coord := models.Coordinator{
		Name:    "coord-0",
		URL:     mustUrl("http://trino.local:8080"),
		Tags:    map[string]string{},
		Enabled: true,
	}

	err := pool.Add(coord)
	require.NoError(t, err)

	err = pool.UpdateStatus()
	require.NoError(t, err)

	backends := pool.Fetch(FetchRequest{})
	require.Len(t, backends, 1)

	availables := pool.Fetch(FetchRequest{
		Status: ClusterStatusEnabled,
		Health: healthcheck.StatusHealthy,
	})

	require.Empty(t, availables)
}

func TestPool_RemoveBackend(t *testing.T) {

	sessStore := session.NewMemoryStorage()
	hc := healthcheck.Mock(healthcheck.Health{
		Status:    healthcheck.StatusHealthy,
		Message:   "ok",
		Timestamp: time.Now(),
	}, nil)

	stats := trino.Mock(models.ClusterStatistics{}, nil)
	logger := logging.Noop()

	pool := NewPool(PoolConfigTest(), sessStore, hc, stats, logger)

	err := pool.Add(models.Coordinator{
		Name:    "coord-0",
		URL:     mustUrl("http://trino.local:8080"),
		Tags:    map[string]string{},
		Enabled: true,
	})
	require.NoError(t, err)

	err = pool.Add(models.Coordinator{
		Name:    "coord-1",
		URL:     mustUrl("http://trino.local:8080"),
		Tags:    map[string]string{},
		Enabled: true,
	})
	require.NoError(t, err)

	// force pool health update before further checking
	// at the moment the health check is executed synchronously in the Add() method but this behaviour may change
	// for the purpose of this test we don't care about the Add method behaviour
	err = pool.UpdateStatus()
	require.NoError(t, err)

	backends := pool.Fetch(FetchRequest{})
	require.Len(t, backends, 2)

	toRemove := pool.Fetch(FetchRequest{
		Name: "coord-0",
	})

	require.Len(t, toRemove, 1)

	err = pool.Remove(toRemove[0].ID)
	require.NoError(t, err)

	backends = pool.Fetch(FetchRequest{})
	require.Len(t, backends, 1)

	coordsByName := pool.Fetch(FetchRequest{
		Name:   "coord-1",
		Health: healthcheck.StatusHealthy,
	})

	require.Len(t, coordsByName, 1)
}

func TestPool_GetByName(t *testing.T) {

	sessStore := session.NewMemoryStorage()
	hc := healthcheck.Mock(healthcheck.Health{
		Status:    healthcheck.StatusHealthy,
		Message:   "ok",
		Timestamp: time.Now(),
	}, nil)

	stats := trino.Mock(models.ClusterStatistics{}, nil)
	logger := logging.Noop()

	pool := NewPool(PoolConfigTest(), sessStore, hc, stats, logger)

	coordsByName := pool.Fetch(FetchRequest{
		Name:   "test",
		Health: healthcheck.StatusHealthy,
	})

	require.Len(t, coordsByName, 0)
}

func TestPool_GetByNameWithUnhealthyStatus(t *testing.T) {
	sessStore := session.NewMemoryStorage()
	hc := healthcheck.Mock(healthcheck.Health{
		Status:    healthcheck.StatusUnhealthy,
		Message:   "ok",
		Timestamp: time.Now(),
	}, nil)

	stats := trino.Mock(models.ClusterStatistics{}, nil)
	logger := logging.Noop()

	pool := NewPool(PoolConfigTest(), sessStore, hc, stats, logger)
	err := pool.Add(models.Coordinator{
		Name:    "coord-1",
		URL:     mustUrl("http://trino.local:8080"),
		Tags:    map[string]string{},
		Enabled: true,
	})
	require.NoError(t, err)

	err = pool.UpdateStatus()
	require.NoError(t, err)
	coordsByName := pool.Fetch(FetchRequest{
		Name:   "coord-1",
		Health: healthcheck.StatusUnhealthy,
	})

	require.Len(t, coordsByName, 1)
}

func TestPool_UpdateBackend(t *testing.T) {
	sessStore := session.NewMemoryStorage()
	hc := healthcheck.Mock(healthcheck.Health{
		Status:    healthcheck.StatusHealthy,
		Message:   "ok",
		Timestamp: time.Now(),
	}, nil)

	stats := trino.Mock(models.ClusterStatistics{}, nil)
	logger := logging.Noop()

	pool := NewPool(PoolConfigTest(), sessStore, hc, stats, logger)

	coord := models.Coordinator{
		Name:    "coord-0",
		URL:     mustUrl("http://trino.local:8080"),
		Tags:    map[string]string{},
		Enabled: true,
	}

	err := pool.Add(coord)
	require.NoError(t, err)

	state := pool.Fetch(FetchRequest{
		Name:   coord.Name,
		Health: healthcheck.StatusHealthy,
	})

	require.Len(t, state, 1)

	require.Equal(t, state[0].Coordinator, coord)

	newState := models.Coordinator{
		Tags: map[string]string{
			"updated": "true",
		},
		Enabled: false,
	}
	err = pool.Update(state[0].ID, newState)
	require.NoError(t, err)

	state = pool.Fetch(FetchRequest{
		Name:   coord.Name,
		Health: healthcheck.StatusHealthy,
	})

	require.Equal(t, models.Coordinator{
		Name:    coord.Name,
		URL:     coord.URL,
		Tags:    newState.Tags,
		Enabled: newState.Enabled,
	}, state[0].Coordinator)

}
