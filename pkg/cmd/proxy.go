package cmd

import (
	"fmt"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/api"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/factory"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/lb"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/ui"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"net/http"
)

var (
	enableProfiling bool
	profilingAddr   string
	staticFilesPath string
)

func init() {
	proxyCmd.Flags().StringVar(&profilingAddr, "profile-addr", ":6060", "profiling server addr")
	proxyCmd.Flags().BoolVar(&enableProfiling, "profile", false, "enable profiling server")
	proxyCmd.Flags().StringVar(&staticFilesPath, "static-files", "ui/dist", "static resources path (ui)")
	rootCmd.AddCommand(proxyCmd)
}

var proxyCmd = &cobra.Command{
	Use:   "proxy",
	Short: "start load balancer",
	Run: func(cmd *cobra.Command, args []string) {

		if enableProfiling {
			go enableProfilingServer(profilingAddr)
		}

		router, err := factory.CreateQueryRouter(factory.QueryRouterConfiguration{
			Type: viper.GetString("routing.rule"),
		})

		if err != nil {
			log.Fatal(err)
		}

		poolConfig := lb.PoolConfig{
			HealthCheckDelay: viper.GetDuration("clusters.healthcheck.delay"),
			StatisticsDelay:  viper.GetDuration("clusters.statistics.delay"),
		}

		pool := lb.NewPool(poolConfig, sessionStorage, clusterHealthCheck, clusterStats, logger)
		sync := lb.NewPoolStateSync(discoveryStorage, logger)

		logger.Info("proxy initialized, syncing cluster state")

		if err := sync.Sync(pool); err != nil {
			log.Fatal(err)
		}

		logger.Info("cluster state sync success")

		conf := lb.ProxyConf{
			SyncDelay: viper.GetDuration("clusters.sync.delay"),
		}

		proxy := lb.NewProxy(conf, pool, sync, sessionStorage, router, logger)
		if err := proxy.Init(); err != nil {
			log.Fatal(err)
		}

		port := viper.GetInt("proxy.port")

		logger.Info("proxy listening on port %d", port)

		httpRouter := mux.NewRouter()

		api := api.NewApi(clusterStats, discover, discoveryStorage, logger)
		uiSrv := ui.New(staticFilesPath)

		httpRouter.PathPrefix("/ui").Handler(uiSrv.Router())
		httpRouter.PathPrefix("/api").Handler(api.Router())
		httpRouter.PathPrefix("/").Handler(proxy.Router())

		corsOpts := cors.AllowAll()

		srv := &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: corsOpts.Handler(httpRouter),
		}

		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	},
}
