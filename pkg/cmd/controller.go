package cmd

import (
	"context"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/api/trino"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/configuration"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/controller/process"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	controllerDelay time.Duration
)

func init() {
	controllerCmd.Flags().DurationVar(&controllerDelay, "every", 10*time.Second, "delay between controller run")
	rootCmd.AddCommand(controllerCmd)
}

var controllerCmd = &cobra.Command{
	Use:   "controller",
	Short: "start trino cluster controller",
	Run: func(cmd *cobra.Command, args []string) {
		var conf configuration.ControllerConf
		if err := viper.UnmarshalKey("controller", &conf); err != nil {
			log.Fatal(err)
		}

		handlers, err := configuration.CreateHandlers(redisClient, logger, notifiers, conf)
		if err != nil {
			log.Fatal(err)
		}

		var ctrl = process.NewController(
			trino.NewClusterApi(),
			discoveryStorage,
			clusterHealthCheck,
			process.NewRedisControllerState(redisClient),
			handlers,
			logger,
		)

		tick := time.NewTicker(controllerDelay)
		ctx, cancel := context.WithCancel(context.Background())

		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		var wg = &sync.WaitGroup{}
		wg.Add(1)

		go func() {
			for {
				select {
				case <-tick.C:
					if err := ctrl.Run(ctx); err != nil {
						log.Fatal(err)
					}
				case <-sigs:
					cancel()
				case <-ctx.Done():
					tick.Stop()
					wg.Done()
					return
				}
			}
		}()

		wg.Wait()
	},
}
