package lb

import (
	"context"
	"database/sql"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/api/trino"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/logging"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/models"

	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/tests"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/discovery"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/proxy/healthcheck"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/proxy/routing"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/proxy/session"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
	"time"
)
import _ "github.com/trinodb/trino-go-client/trino"

func TestIntegration(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	cluster0, c0, err := tests.CreateTrinoCluster(ctx)
	require.NoError(t, err)
	defer func() {
		// if the test ran successfully this container should be already terminated
		_ = cluster0.Terminate(ctx)
	}()

	cluster1, c1, err := tests.CreateTrinoCluster(ctx)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, cluster1.Terminate(ctx))
	}()

	cluster2, c2, err := tests.CreateTrinoCluster(ctx)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, cluster2.Terminate(ctx))
	}()

	stateStore := discovery.NewMemoryStorage()
	sessStore := session.NewMemoryStorage()
	hc := healthcheck.NewHttpHealth()
	stats := trino.NewClusterApi()

	router := routing.New(routing.NewUserAwareRouter(routing.UserAwareRoutingConf{}), routing.RoundRobin())

	logger := logging.Noop()

	poolCfg := PoolConfig{
		HealthCheckDelay: 5 * time.Second,
		StatisticsDelay:  5 * time.Second,
	}

	pool := NewPool(poolCfg, sessStore, hc, stats, logger)
	poolSync := NewPoolStateSync(stateStore, logging.Noop())

	err = stateStore.Add(ctx, models.Coordinator{
		Name:    "c0",
		URL:     c0,
		Tags:    nil,
		Enabled: true,
	})
	require.NoError(t, err)

	err = stateStore.Add(ctx, models.Coordinator{
		Name:    "c1",
		URL:     c1,
		Tags:    nil,
		Enabled: true,
	})
	require.NoError(t, err)

	err = stateStore.Add(ctx, models.Coordinator{
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
