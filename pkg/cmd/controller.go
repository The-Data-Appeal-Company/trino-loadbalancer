package cmd

import (
	"github.com/spf13/cobra"
	"time"
)

var (
	controllerDelay time.Duration
)

func init() {
	controllerCmd.PersistentFlags().DurationVar(&controllerDelay, "every", 10*time.Second, "delay between controller run")
	rootCmd.AddCommand(controllerCmd)
}

var controllerCmd = &cobra.Command{
	Use:   "controller",
	Short: "start trino cluster controller [legacy, autoscaler]",
}
