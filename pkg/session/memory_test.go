package session

import (
	"context"
	"errors"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/models"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMemoryLinkerLinkCluster(t *testing.T) {
	t.Parallel()

	storage := NewMemoryStorage()
	ctx := context.TODO()

	queryInfo := models.QueryInfo{
		User:          "user",
		QueryID:       "query",
		TransactionID: "tx",
	}

	const coordinator = "benny-lo-spenny"
	err := storage.Link(ctx, queryInfo, coordinator)
	require.NoError(t, err)

	linkedCoordinator, err := storage.Get(ctx, queryInfo)

	require.NoError(t, err)
	require.Equal(t, coordinator, linkedCoordinator)

	err = storage.Unlink(ctx, queryInfo)
	require.NoError(t, err)

	_, err = storage.Get(ctx, queryInfo)

	require.Error(t, err)
	require.True(t, errors.Is(err, ErrLinkNotFound))

}

func TestMemoryLinkerLinkNotFoundErr(t *testing.T) {
	t.Parallel()

	storage := NewMemoryStorage()
	ctx := context.TODO()

	_, err := storage.Get(ctx, models.QueryInfo{
		User:          "user",
		QueryID:       "query",
		TransactionID: "tx",
	})

	require.Error(t, err)
	require.True(t, errors.Is(err, ErrLinkNotFound))

}
