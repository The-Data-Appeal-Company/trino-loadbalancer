package lb

import (
	"context"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/discovery"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/healthcheck"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/logging"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/models"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/tests"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSyncPoolStatus_AddMissingFromState(t *testing.T) {
	ctx := context.Background()

	pool := NewMockPool()
	storage := discovery.NewMemoryStorage()

	coord := models.Coordinator{
		Name:    "test-0",
		URL:     tests.MustUrl("http://presto.local:8889"),
		Enabled: true,
		Tags: map[string]string{
			"test": "asd",
		},
	}

	err := storage.Add(ctx, coord)

	require.NoError(t, err)

	sync := NewPoolStateSync(storage, logging.Noop())

	err = sync.Sync(pool)
	require.NoError(t, err)

	backends := pool.AllBackends()
	require.Len(t, backends, 1)
	require.Equal(t, backends[0].Backend.Name, coord.Name)
	require.Equal(t, backends[0].Backend.URL.String(), coord.URL.String())
	require.Equal(t, backends[0].Backend.Enabled, coord.Enabled)

}

func TestSyncPoolStatus_Remove(t *testing.T) {
	pool := NewMockPool()
	storage := discovery.NewMemoryStorage()

	coord := models.Coordinator{
		Name:    "test-0",
		URL:     tests.MustUrl("http://presto.local:8889"),
		Enabled: true,
		Tags: map[string]string{
			"test": "asd",
		},
	}

	err := pool.Add(coord)
	require.NoError(t, err)

	sync := NewPoolStateSync(storage, logging.Noop())

	err = sync.Sync(pool)
	require.NoError(t, err)

	backends := pool.AllBackends()
	require.Len(t, backends, 0)

}

func TestSyncPoolStatus_DoNothing(t *testing.T) {
	ctx := context.Background()

	pool := NewMockPool()
	storage := discovery.NewMemoryStorage()

	coord0 := models.Coordinator{
		Name:    "test-0",
		URL:     tests.MustUrl("http://presto0.local:8889"),
		Enabled: true,
		Tags: map[string]string{
			"test": "asd",
		},
	}

	coord1 := models.Coordinator{
		Name:    "test-1",
		URL:     tests.MustUrl("http://presto1.local:8889"),
		Enabled: true,
		Tags: map[string]string{
			"test": "asd",
		},
	}

	err := pool.Add(coord0)
	require.NoError(t, err)

	err = pool.Add(coord1)
	require.NoError(t, err)

	err = storage.Add(ctx, coord1)
	require.NoError(t, err)

	err = storage.Add(ctx, coord0)
	require.NoError(t, err)

	sync := NewPoolStateSync(storage, logging.Noop())

	for i := 0; i < 3; i++ {
		err = sync.Sync(pool)
		require.NoError(t, err)
	}

	stateBackends, err := storage.All(ctx)
	require.NoError(t, err)

	require.Len(t, pool.AllBackends(), len(stateBackends))

}

func TestSyncPoolStatus_UpdateEnabledStatus(t *testing.T) {
	ctx := context.Background()

	pool := NewMockPool()
	storage := discovery.NewMemoryStorage()

	coord0 := models.Coordinator{
		Name:    "test-0",
		URL:     tests.MustUrl("http://presto0.local:8889"),
		Enabled: true,
		Tags: map[string]string{
			"test": "asd",
		},
	}

	coord1 := models.Coordinator{
		Name:    "test-1",
		URL:     tests.MustUrl("http://presto1.local:8889"),
		Enabled: true,
		Tags: map[string]string{
			"test": "asd",
		},
	}

	err := pool.Add(coord0)
	require.NoError(t, err)

	err = pool.Add(coord1)
	require.NoError(t, err)

	err = storage.Add(ctx, coord1)
	require.NoError(t, err)

	err = storage.Add(ctx, coord0)
	require.NoError(t, err)

	sync := NewPoolStateSync(storage, logging.Noop())

	err = sync.Sync(pool)
	require.NoError(t, err)

	// Simulate coord0 status change
	err = storage.Remove(ctx, coord0.Name)
	require.NoError(t, err)

	coord0.Enabled = false
	err = storage.Add(ctx, coord0)
	require.NoError(t, err)

	// sync after status change
	err = sync.Sync(pool)
	require.NoError(t, err)

	coord0FromPool, err := pool.GetByName(coord0.Name, healthcheck.StatusHealthy)
	require.NoError(t, err)

	require.Equal(t, coord0FromPool.Backend, coord0)

}
