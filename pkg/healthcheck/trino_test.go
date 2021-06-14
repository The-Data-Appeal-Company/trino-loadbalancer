package healthcheck

import (
	"context"
	"fmt"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/tests"
	"github.com/stretchr/testify/require"
	"net/url"
	"testing"
)

func TestTrinoClusterHealth_Check(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	container, _, err := tests.CreateTrinoCluster(ctx)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, container.Terminate(ctx))
	}()

	check := NewHttpHealth()

	port, err := container.MappedPort(ctx, "8080")
	require.NoError(t, err)

	uri, err := url.Parse(fmt.Sprintf("http://localhost:%d", port.Int()))
	require.NoError(t, err)

	result, err := check.Check(uri)
	require.NoError(t, err)

	require.True(t, result.IsAvailable(), result.Message)
}

func TestTrinoClusterHealth_CheckDown(t *testing.T) {
	check := NewHttpHealth()

	uri, err := url.Parse("http://trino.local:8080")
	require.NoError(t, err)

	result, err := check.Check(uri)
	require.NoError(t, err)

	require.False(t, result.IsAvailable(), result.Message)
}
