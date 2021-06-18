package cmd

import (
	"context"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/api/trino"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/configuration"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/controller/process"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"time"
)

func init() {
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

		handlers, err := configuration.CreateHandlers(redisClient, logger, conf)
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

		for {
			if err := ctrl.Run(context.Background()); err != nil {
				log.Fatal(err)
			}
			time.Sleep(1 * time.Second)
		}
	},
}
