package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"time"
)

var (
	discoveryDelay time.Duration
)

func init() {
	discoverCmd.Flags().DurationVar(&discoveryDelay, "every", 10*time.Minute, "delay between discovery run")
	rootCmd.AddCommand(discoverCmd)
}

var discoverCmd = &cobra.Command{
	Use:   "discovery",
	Short: "start presto cluster discovery",
	Run: func(cmd *cobra.Command, args []string) {

		for {
			clusters, err := discover.Discover()
			if err != nil {
				log.Fatal(err)
			}

			for _, cluster := range clusters {
				if err := discoveryStorage.Add(cmd.Context(), cluster); err != nil {
					log.Fatal(err)
				}
				logger.Info("found cluster: %s ( %s )", cluster.Name, cluster.URL)
			}

			time.Sleep(discoveryDelay)
		}
	},
}
