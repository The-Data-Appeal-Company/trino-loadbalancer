package lb

import (
	"context"
	logging2 "github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/logging"
	models2 "github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/models"
	tests2 "github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/tests"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/discovery"
	healthcheck2 "github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/proxy/healthcheck"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSyncPoolStatus_AddMissingFromState(t *testing.T) {
	ctx := context.Background()

	pool := NewMockPool()
	storage := discovery.NewMemoryStorage()

	coord := models2.Coordinator{
		Name:    "test-0",
		URL:     tests2.MustUrl("http://trino.local:8889"),
		Enabled: true,
		Tags: map[string]string{
			"test": "asd",
		},
	}

	err := storage.Add(ctx, coord)

	require.NoError(t, err)

	sync := NewPoolStateSync(storage, logging2.Noop())

	err = sync.Sync(pool)
	require.NoError(t, err)

	backends := pool.Fetch(FetchRequest{})
	require.Len(t, backends, 1)
	require.Equal(t, backends[0].Coordinator.Name, coord.Name)
	require.Equal(t, backends[0].Coordinator.URL.String(), coord.URL.String())
	require.Equal(t, backends[0].Coordinator.Enabled, coord.Enabled)

}

func TestSyncPoolStatus_Remove(t *testing.T) {
	pool := NewMockPool()
	storage := discovery.NewMemoryStorage()

	coord := models2.Coordinator{
		Name:    "test-0",
		URL:     tests2.MustUrl("http://trino.local:8889"),
		Enabled: true,
		Tags: map[string]string{
			"test": "asd",
		},
	}

	err := pool.Add(coord)
	require.NoError(t, err)

	sync := NewPoolStateSync(storage, logging2.Noop())

	err = sync.Sync(pool)
	require.NoError(t, err)

	backends := pool.Fetch(FetchRequest{})
	require.Len(t, backends, 0)

}

func TestSyncPoolStatus_DoNothing(t *testing.T) {
	ctx := context.Background()

	pool := NewMockPool()
	storage := discovery.NewMemoryStorage()

	coord0 := models2.Coordinator{
		Name:    "test-0",
		URL:     tests2.MustUrl("http://trino0.local:8889"),
		Enabled: true,
		Tags: map[string]string{
			"test": "asd",
		},
	}

	coord1 := models2.Coordinator{
		Name:    "test-1",
		URL:     tests2.MustUrl("http://trino1.local:8889"),
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

	sync := NewPoolStateSync(storage, logging2.Noop())

	for i := 0; i < 3; i++ {
		err = sync.Sync(pool)
		require.NoError(t, err)
	}

	stateBackends, err := storage.All(ctx)
	require.NoError(t, err)

	require.Len(t, pool.Fetch(FetchRequest{}), len(stateBackends))

}

func TestSyncPoolStatus_UpdateEnabledStatus(t *testing.T) {
	ctx := context.Background()

	pool := NewMockPool()
	storage := discovery.NewMemoryStorage()

	coord0 := models2.Coordinator{
		Name:    "test-0",
		URL:     tests2.MustUrl("http://trino0.local:8889"),
		Enabled: true,
		Tags: map[string]string{
			"test": "asd",
		},
	}

	coord1 := models2.Coordinator{
		Name:    "test-1",
		URL:     tests2.MustUrl("http://trino1.local:8889"),
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

	sync := NewPoolStateSync(storage, logging2.Noop())

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

	coord0FromPool := pool.Fetch(FetchRequest{
		Name:   coord0.Name,
		Health: healthcheck2.StatusHealthy,
	})
	require.NoError(t, err)

	require.Equal(t, coord0FromPool[0].Coordinator, coord0)

}
