package session

import (
	"context"
	"errors"
	models2 "github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/models"
	tests2 "github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/tests"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestRedisLinkerLinkCluster(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	container, client, err := tests2.CreateRedisServer(ctx)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, container.Terminate(ctx))
	}()

	storage := RedisLinkerStorage{
		redis:  client,
		prefix: "test",
		ttl:    10 * time.Minute,
	}

	queryInfo := models2.QueryInfo{
		User:          "user",
		QueryID:       "query",
		TransactionID: "tx",
	}

	const coordinator = "benny-lo-spenny"
	err = storage.Link(ctx, queryInfo, coordinator)
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

func TestRedisLinkerLinkNotFoundErr(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	container, client, err := tests2.CreateRedisServer(ctx)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, container.Terminate(ctx))
	}()

	storage := RedisLinkerStorage{
		redis:  client,
		prefix: "test",
		ttl:    10 * time.Minute,
	}

	_, err = storage.Get(ctx, models2.QueryInfo{
		User:          "user",
		QueryID:       "query",
		TransactionID: "tx",
	})

	require.Error(t, err)
	require.True(t, errors.Is(err, ErrLinkNotFound))

}
