package healthcheck

import (
	"context"
	"fmt"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/tests"
	"github.com/stretchr/testify/require"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
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

func TestHttpClient_TimeoutWithConnHang(t *testing.T) {
	backendSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(1 * time.Second)
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte("Internal Server Error"))
		require.NoError(t, err)
	}))
	defer backendSrv.Close()

	uri, err := url.Parse(backendSrv.URL)
	require.NoError(t, err)
	check := NewHttpHealthWithTimeout(300 * time.Millisecond)

	_, err = check.client.Get(uri.String())

	netErr, isNetErr := err.(net.Error)
	require.Error(t, netErr)
	require.True(t, isNetErr)
}

func TestTrinoClusterHealth_Check_FailOnConnectionHang(t *testing.T) {
	backendSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(1000 * time.Millisecond)
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte("Internal Server Error"))
		require.NoError(t, err)

	}))
	defer backendSrv.Close()

	uri, err := url.Parse(backendSrv.URL)
	require.NoError(t, err)

	check := NewHttpHealthWithTimeout(10 * time.Millisecond)

	result, err := check.Check(uri)
	require.NoError(t, err)
	require.False(t, result.IsAvailable())
	backendSrv.Close()
}
