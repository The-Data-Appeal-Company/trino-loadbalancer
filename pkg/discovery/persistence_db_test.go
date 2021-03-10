package discovery

import (
	"context"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/models"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/tests"
	"github.com/stretchr/testify/require"
	"testing"
)

// To avoid spinning up multiple postgres containers there are as few tests are possible
// we should use some library that allow setUp / tearDown for test suites
func TestDbPersistence(t *testing.T) {

	ctx := context.Background()
	container, db, err := tests.CreatePostgresDatabase(ctx, tests.WithInitScript("testdata/init.sql"))
	require.NoError(t, err)

	defer container.Terminate(ctx)

	storage := NewDatabaseStorage(db, DefaultDatabaseTableName)

	coord0 := models.Coordinator{
		Name: "test-0",
		URL:  tests.MustUrl("http://test.local:8889"),
		Tags: map[string]string{
			"test": "true",
		},
		Enabled:      true,
	}
	err = storage.Add(ctx, coord0)
	require.NoError(t, err)

	backends, err := storage.All(ctx)
	require.NoError(t, err)
	require.Len(t, backends, 1)
	require.Equal(t, coord0.Name, backends[0].Name)
	require.Equal(t, coord0.URL.String(), backends[0].URL.String())
	require.Equal(t, coord0.Enabled, backends[0].Enabled)
	require.Equal(t, coord0.Tags, backends[0].Tags)

	backend, err := storage.Get(ctx, coord0.Name)
	require.NoError(t, err)

	require.Equal(t, coord0.Name, backend.Name)
	require.Equal(t, coord0.URL.String(), backend.URL.String())
	require.Equal(t, coord0.Enabled, backend.Enabled)
	require.Equal(t, coord0.Tags, backend.Tags)

	err = storage.Remove(ctx, coord0.Name)
	require.NoError(t, err)

	all, err := storage.All(ctx)
	require.NoError(t, err)
	require.Empty(t, all)

	_, err = storage.Get(ctx, coord0.Name)
	require.Error(t, err)
	require.Equal(t, err, ErrClusterNotFound)

	// return error when adding multiple clusters with same name
	err = storage.Add(ctx, coord0)
	require.NoError(t, err)
	err = storage.Add(ctx, coord0)
	require.NoError(t, err)

}

func TestDbDoubleInsertUpdate(t *testing.T) {

	ctx := context.Background()
	container, db, err := tests.CreatePostgresDatabase(ctx, tests.WithInitScript("testdata/init.sql"))
	require.NoError(t, err)

	defer container.Terminate(ctx)

	storage := NewDatabaseStorage(db, DefaultDatabaseTableName)

	err = storage.Add(ctx, models.Coordinator{
		Name: "test-0",
		URL:  tests.MustUrl("http://test.local:8889"),
		Tags: map[string]string{
			"test": "true",
		},
		Enabled: true,
	})

	require.NoError(t, err)

	updated := models.Coordinator{
		Name: "test-0",
		URL:  tests.MustUrl("http://test.local:8889"),
		Tags: map[string]string{
			"test":    "true",
			"updated": "true",
		},
		Enabled: false,
	}
	err = storage.Add(ctx, updated)
	require.NoError(t, err)

	all, err := storage.All(ctx)
	require.NoError(t, err)

	require.Len(t, all, 1)
	require.Equal(t, all[0], updated)
}
