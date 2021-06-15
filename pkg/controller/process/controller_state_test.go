package process

import (
	"context"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/models"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/tests"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestSetAndGetState(t *testing.T) {

	ctx := context.Background()
	container, redis, err := tests.CreateRedisServer(ctx)
	defer func() {
		require.NoError(t, container.Terminate(ctx))
	}()
	require.NoError(t, err)
	state := NewRedisControllerState(redis)

	firstCluster := models.Coordinator{
		Name: "node-00",
	}
	secondCluster := models.Coordinator{
		Name: "node-01",
	}

	noTime, err := state.Get(ctx, firstCluster)
	require.NoError(t, err)
	require.True(t, noTime.Equal(time.Unix(0, 0)))

	processTime := time.Date(2020, 3, 8, 0, 0, 0, 0, time.UTC)
	require.NoError(t, state.Set(ctx, firstCluster, processTime))

	retrievedTime, err := state.Get(ctx, firstCluster)
	require.NoError(t, err)
	require.True(t, retrievedTime.Equal(processTime))

	noTime, err = state.Get(ctx, secondCluster)
	require.NoError(t, err)
	require.True(t, noTime.Equal(time.Unix(0, 0)))
}
