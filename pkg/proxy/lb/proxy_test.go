package lb

import (
	trino2 "github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/api/trino"
	logging2 "github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/logging"
	models2 "github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/models"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/discovery"
	healthcheck2 "github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/proxy/healthcheck"
	routing2 "github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/proxy/routing"
	session2 "github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/proxy/session"
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

	sessStore := session2.NewMemoryStorage()
	hc := healthcheck2.NoOp()
	stats := trino2.Noop()

	router := routing2.New(routing2.NewUserAwareRouter(routing2.UserAwareRoutingConf{}), routing2.RandomRouter{})

	logger := logging2.Noop()
	pool := NewPool(PoolConfigTest(), sessStore, hc, stats, logger)

	err := pool.Add(models2.Coordinator{
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

	sessStore := session2.NewMemoryStorage()
	hc := healthcheck2.NoOp()
	stats := trino2.Noop()

	router := routing2.New(routing2.NewUserAwareRouter(routing2.UserAwareRoutingConf{}), routing2.RoundRobin())

	logger := logging2.Noop()
	pool := NewPool(PoolConfigTest(), sessStore, hc, stats, logger)

	err := pool.Add(models2.Coordinator{
		Name:    "cluster-0",
		URL:     mustUrl(fakeCoord0.URL),
		Enabled: true,
	})
	require.NoError(t, err)

	err = pool.Add(models2.Coordinator{
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

	sessStore := session2.NewMemoryStorage()
	hc := healthcheck2.NewHttpHealth()
	stats := trino2.Noop()

	router := routing2.New(routing2.NewUserAwareRouter(routing2.UserAwareRoutingConf{}), routing2.RandomRouter{})

	logger := logging2.Noop()

	pool := NewPool(PoolConfigTest(), sessStore, hc, stats, logger)
	err := pool.Add(models2.Coordinator{
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

	sessStore := session2.NewMemoryStorage()
	hc := healthcheck2.NoOp()
	stats := trino2.Noop()

	router := routing2.New(routing2.NewUserAwareRouter(routing2.UserAwareRoutingConf{}), routing2.RandomRouter{})

	logger := logging2.Noop()
	pool := NewPool(PoolConfigTest(), sessStore, hc, stats, logger)

	err := pool.Add(models2.Coordinator{
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
	sessStore := session2.NewMemoryStorage()
	hc := healthcheck2.NewHttpHealth()
	stats := trino2.Noop()

	router := routing2.New(routing2.NewUserAwareRouter(routing2.UserAwareRoutingConf{}), routing2.RandomRouter{})

	logger := logging2.Noop()

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
	sessStore := session2.NewMemoryStorage()
	hc := healthcheck2.NewHttpHealth()
	stats := trino2.Noop()

	router := routing2.New(routing2.NewUserAwareRouter(routing2.UserAwareRoutingConf{}), routing2.RandomRouter{})

	logger := logging2.Noop()

	pool := NewPool(PoolConfigTest(), sessStore, hc, stats, logger)
	poolStateSync := NewPoolStateSync(stateStore, logger)

	proxy := NewProxy(proxyConfig, pool, poolStateSync, sessStore, router, logger)

	srv := httptest.NewServer(proxy.Router())
	defer srv.Close()

	res, err := http.Get(srv.URL + "/health")
	require.NoError(t, err)
	require.Equal(t, res.StatusCode, http.StatusOK)

}
