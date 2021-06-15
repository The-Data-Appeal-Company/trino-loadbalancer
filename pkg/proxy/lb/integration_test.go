package lb

import (
	"context"
	"database/sql"
	trino2 "github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/api/trino"
	logging2 "github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/logging"
	models2 "github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/models"
	tests2 "github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/tests"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/discovery"
	healthcheck2 "github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/proxy/healthcheck"
	routing2 "github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/proxy/routing"
	session2 "github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/proxy/session"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
	"time"
)
import _ "github.com/trinodb/trino-go-client/trino"

func TestIntegration(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	cluster0, c0, err := tests2.CreateTrinoCluster(ctx)
	require.NoError(t, err)
	defer func() {
		// if the test ran successfully this container should be already terminated
		_ = cluster0.Terminate(ctx)
	}()

	cluster1, c1, err := tests2.CreateTrinoCluster(ctx)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, cluster1.Terminate(ctx))
	}()

	cluster2, c2, err := tests2.CreateTrinoCluster(ctx)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, cluster2.Terminate(ctx))
	}()

	stateStore := discovery.NewMemoryStorage()
	sessStore := session2.NewMemoryStorage()
	hc := healthcheck2.NewHttpHealth()
	stats := trino2.NewClusterApi()

	router := routing2.New(routing2.NewUserAwareRouter(routing2.UserAwareRoutingConf{}), routing2.RoundRobin())

	logger := logging2.Noop()

	poolCfg := PoolConfig{
		HealthCheckDelay: 5 * time.Second,
		StatisticsDelay:  5 * time.Second,
	}

	pool := NewPool(poolCfg, sessStore, hc, stats, logger)
	poolSync := NewPoolStateSync(stateStore, logging2.Noop())

	err = stateStore.Add(ctx, models2.Coordinator{
		Name:    "c0",
		URL:     c0,
		Tags:    nil,
		Enabled: true,
	})
	require.NoError(t, err)

	err = stateStore.Add(ctx, models2.Coordinator{
		Name:    "c1",
		URL:     c1,
		Tags:    nil,
		Enabled: true,
	})
	require.NoError(t, err)

	err = stateStore.Add(ctx, models2.Coordinator{
		Name:    "c2",
		URL:     c2,
		Tags:    nil,
		Enabled: true,
	})
	require.NoError(t, err)

	err = poolSync.Sync(pool)
	require.NoError(t, err)

	proxy := NewProxy(proxyConfig, pool, poolSync, sessStore, router, logger)

	go func() {
		require.NoError(t, proxy.Init())
		require.NoError(t, proxy.Serve("0.0.0.0:4322"))
	}()

	time.Sleep(300 * time.Millisecond)

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			err = testQuery("http://test@localhost:4322?catalog=memory")
			require.NoError(t, err)
		}(&wg)
	}
	wg.Wait()

	err = cluster0.Terminate(ctx)
	require.NoError(t, err)

	err = pool.UpdateStatus()
	require.NoError(t, err)

	for i := 0; i < 5; i++ {
		err = testQuery("http://test@localhost:4322?catalog=memory&schema=test")
		require.NoError(t, err)
	}

}

func testQuery(address string) error {
	db, err := sql.Open("trino", address)
	if err != nil {
		return err
	}

	row, err := db.Query("select 1")
	if err != nil {
		return err
	}

	row.Next()

	var r int
	if err := row.Scan(&r); err != nil {
		return err
	}

	return nil
}
