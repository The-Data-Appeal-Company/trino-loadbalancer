package discovery

import (
	"context"
	models2 "github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/models"
	tests2 "github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/tests"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMemoryStorage_Add(t *testing.T) {
	ctx := context.TODO()
	storage := NewMemoryStorage()

	coord0 := models2.Coordinator{
		Name:    "coord-0",
		URL:     tests2.MustUrl("http://trino.local:8080"),
		Tags:    map[string]string{},
		Enabled: true,
	}
	err := storage.Add(ctx, coord0)
	require.NoError(t, err)

	backends, err := storage.All(ctx)
	require.NoError(t, err)
	require.Len(t, backends, 1)
	require.Equal(t, coord0.Name, backends[0].Name)
	require.Equal(t, coord0.URL.String(), backends[0].URL.String())
	require.Equal(t, coord0.Enabled, backends[0].Enabled)
	require.Equal(t, coord0.Tags, backends[0].Tags)
}

func TestMemoryStorage_AddRemove(t *testing.T) {
	ctx := context.TODO()
	storage := NewMemoryStorage()

	coord0 := models2.Coordinator{
		Name:    "coord-0",
		URL:     tests2.MustUrl("http://trino.local:8080"),
		Tags:    map[string]string{},
		Enabled: true,
	}
	err := storage.Add(ctx, coord0)
	require.NoError(t, err)

	all, err := storage.All(ctx)
	require.NoError(t, err)
	require.Len(t, all, 1)

	err = storage.Remove(ctx, coord0.Name)
	require.NoError(t, err)

	all, err = storage.All(ctx)
	require.NoError(t, err)
	require.Len(t, all, 0)

}

func TestMemoryStorage_AddGet(t *testing.T) {
	ctx := context.TODO()
	storage := NewMemoryStorage()

	coord0 := models2.Coordinator{
		Name:    "coord-0",
		URL:     tests2.MustUrl("http://trino.local:8080"),
		Tags:    map[string]string{},
		Enabled: true,
	}
	err := storage.Add(ctx, coord0)
	require.NoError(t, err)

	coord, err := storage.Get(ctx, coord0.Name)
	require.NoError(t, err)
	require.Equal(t, coord0.Name, coord.Name)
	require.Equal(t, coord0.URL.String(), coord.URL.String())
	require.Equal(t, coord0.Enabled, coord.Enabled)
	require.Equal(t, coord0.Tags, coord.Tags)

}
