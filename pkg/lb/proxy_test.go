package lb

import (
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/discovery"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/healthcheck"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/logging"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/models"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/routing"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/session"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/statistics"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"time"

	"testing"
)

var proxyConfig = ProxyConf{
	SyncDelay: 1 * time.Hour,
}

func TestProxyRouting(t *testing.T) {

	fakeCoord0 := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
	}))
	defer fakeCoord0.Close()

	sessStore := session.NewMemoryStorage()
	hc := healthcheck.NoOp()
	stats := statistics.Noop()

	router := routing.New(routing.RandomRouter{})

	logger := logging.Noop()
	pool := NewPool(PoolConfigTest(), sessStore, hc, stats, logger)

	err := pool.Add(models.Coordinator{
		Name:    "test",
		URL:     mustUrl(fakeCoord0.URL),
		Enabled: true,
	})
	require.NoError(t, err)

	proxy := NewProxy(proxyConfig, pool, NoOpSync{}, sessStore, router, logger)

	srv := httptest.NewServer(http.HandlerFunc(proxy.Handle))
	defer srv.Close()

	res, err := http.Post(srv.URL, "application/json", nil)
	require.NoError(t, err)
	require.Equal(t, res.StatusCode, http.StatusOK)

}

func TestProxyRoutingMultiCoordinator(t *testing.T) {

	fakeCoord0 := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
	}))
	defer fakeCoord0.Close()

	fakeCoord1 := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
	}))
	defer fakeCoord1.Close()

	sessStore := session.NewMemoryStorage()
	hc := healthcheck.NoOp()
	stats := statistics.Noop()

	router := routing.New(routing.RoundRobin())

	logger := logging.Noop()
	pool := NewPool(PoolConfigTest(), sessStore, hc, stats, logger)

	err := pool.Add(models.Coordinator{
		Name:    "cluster-0",
		URL:     mustUrl(fakeCoord0.URL),
		Enabled: true,
	})
	require.NoError(t, err)

	err = pool.Add(models.Coordinator{
		Name:    "cluster-1",
		URL:     mustUrl(fakeCoord1.URL),
		Enabled: true,
	})
	require.NoError(t, err)

	proxy := NewProxy(proxyConfig, pool, NoOpSync{}, sessStore, router, logger)

	srv := httptest.NewServer(http.HandlerFunc(proxy.Handle))
	defer srv.Close()

	for i := 0; i < 10; i++ {
		res, err := http.Post(srv.URL, "application/json", nil)
		require.NoError(t, err)
		require.Equal(t, res.StatusCode, http.StatusOK)
	}
}

func TestProxyWithUnhealthyBackend(t *testing.T) {

	sessStore := session.NewMemoryStorage()
	hc := healthcheck.NewHttpHealth()
	stats := statistics.Noop()

	router := routing.New(routing.RandomRouter{})

	logger := logging.Noop()

	pool := NewPool(PoolConfigTest(), sessStore, hc, stats, logger)
	err := pool.Add(models.Coordinator{
		Name:    "test",
		URL:     mustUrl("http://trino.local:1231"),
		Enabled: true,
	})
	require.NoError(t, err)

	proxy := NewProxy(proxyConfig, pool, NoOpSync{}, sessStore, router, logger)

	srv := httptest.NewServer(http.HandlerFunc(proxy.Handle))
	defer srv.Close()

	res, err := http.Post(srv.URL, "application/json", nil)
	require.NoError(t, err)
	require.Equal(t, res.StatusCode, http.StatusServiceUnavailable)

}

func TestProxyWithHealthyUnreachableBackend(t *testing.T) {

	sessStore := session.NewMemoryStorage()
	hc := healthcheck.NoOp()
	stats := statistics.Noop()

	router := routing.New(routing.RandomRouter{})

	logger := logging.Noop()
	pool := NewPool(PoolConfigTest(), sessStore, hc, stats, logger)

	err := pool.Add(models.Coordinator{
		Name:    "test",
		URL:     mustUrl("http://trino.local:1231"),
		Enabled: true,
	})
	require.NoError(t, err)

	proxy := NewProxy(proxyConfig, pool, NoOpSync{}, sessStore, router, logger)

	srv := httptest.NewServer(http.HandlerFunc(proxy.Handle))
	defer srv.Close()

	res, err := http.Post(srv.URL, "application/json", nil)
	require.NoError(t, err)
	require.Equal(t, res.StatusCode, http.StatusBadGateway)

}

func TestProxyWithNoBackends(t *testing.T) {

	stateStore := discovery.NewMemoryStorage()
	sessStore := session.NewMemoryStorage()
	hc := healthcheck.NewHttpHealth()
	stats := statistics.Noop()

	router := routing.New(routing.RandomRouter{})

	logger := logging.Noop()

	pool := NewPool(PoolConfigTest(), sessStore, hc, stats, logger)
	poolStateSync := NewPoolStateSync(stateStore, logger)

	proxy := NewProxy(proxyConfig, pool, poolStateSync, sessStore, router, logger)

	srv := httptest.NewServer(http.HandlerFunc(proxy.Handle))
	defer srv.Close()

	res, err := http.Post(srv.URL, "application/json", nil)
	require.NoError(t, err)
	require.Equal(t, res.StatusCode, http.StatusServiceUnavailable)

}

func TestProxyHealthEndpoint(t *testing.T) {

	stateStore := discovery.NewMemoryStorage()
	sessStore := session.NewMemoryStorage()
	hc := healthcheck.NewHttpHealth()
	stats := statistics.Noop()

	router := routing.New(routing.RandomRouter{})

	logger := logging.Noop()

	pool := NewPool(PoolConfigTest(), sessStore, hc, stats, logger)
	poolStateSync := NewPoolStateSync(stateStore, logger)

	proxy := NewProxy(proxyConfig, pool, poolStateSync, sessStore, router, logger)

	srv := httptest.NewServer(proxy.Router())
	defer srv.Close()

	res, err := http.Get(srv.URL + "/health")
	require.NoError(t, err)
	require.Equal(t, res.StatusCode, http.StatusOK)

}
