package cmd

import (
	"context"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/api/trino"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/controller/process"
	"github.com/spf13/cobra"
	"log"
)

func init() {
	rootCmd.AddCommand(controllerCmd)
}

var controllerCmd = &cobra.Command{
	Use:   "controller",
	Short: "start trino cluster controller",
	Run: func(cmd *cobra.Command, args []string) {
		var ctrl = process.NewController(trino.NewClusterApi(), discoveryStorage, clusterHealthCheck)
		if err := ctrl.Run(context.Background()); err != nil {
			log.Fatal(err)
		}
	},
}
