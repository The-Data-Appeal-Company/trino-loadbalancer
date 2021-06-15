package session

import (
	"context"
	"errors"
	models2 "github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/models"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMemoryCachedLinkerLinkCluster(t *testing.T) {
	t.Parallel()

	storage := NewCaching(NewMemoryStorage(), NewMemoryStorage())
	ctx := context.TODO()

	queryInfo := models2.QueryInfo{
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

func TestMemoryCachedLinkerLinkNotFoundErr(t *testing.T) {
	t.Parallel()

	storage := NewCaching(NewMemoryStorage(), NewMemoryStorage())
	ctx := context.TODO()

	_, err := storage.Get(ctx, models2.QueryInfo{
		User:          "user",
		QueryID:       "query",
		TransactionID: "tx",
	})

	require.Error(t, err)
	require.True(t, errors.Is(err, ErrLinkNotFound))

}
