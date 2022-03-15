package cmd

import (
	"context"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/api/trino"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/configuration"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/controller/autoscaler"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func init() {
	controllerCmd.AddCommand(autoscalerCmd)
}

var autoscalerCmd = &cobra.Command{
	Use:   "autoscaler",
	Short: "start trino cluster autoscale controller",
	Run: func(cmd *cobra.Command, args []string) {
		var conf configuration.AutoscalerConf
		if err := viper.UnmarshalKey("controller.autoscaler", &conf); err != nil {
			log.Fatal(err)
		}

		tick := time.NewTicker(controllerDelay)
		ctx, cancel := context.WithCancel(context.Background())

		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		// run in-cluster mode
		k8sClient, err := configuration.NewK8sClient(nil)
		if err != nil {
			log.Fatal(err)
		}

		autoscalerController := autoscaler.NewKubeClientAutoscaler(k8sClient, trino.NewClusterApi(), autoscaler.MemoryState(), logger)

		var wg = &sync.WaitGroup{}
		wg.Add(1)

		go func() {
			for {
				select {
				case <-tick.C:

					for _, cluster := range conf.Kubernetes {
						coordUri, err := url.Parse(cluster.CoordinatorUri)
						if err != nil {
							log.Fatal(err)
						}

						err = autoscalerController.Execute(autoscaler.KubeRequest{
							Coordinator:     coordUri,
							Namespace:       cluster.Namespace,
							Deployment:      cluster.Deployment,
							Min:             cluster.Min,
							ScaleUpStrategy: cluster.ScaleUpStrategy,
							ScaleAfter:      cluster.ScaleAfter,
						})

						if err != nil {
							logger.Warn("unable to run autoscaler on %s: %s", cluster.CoordinatorUri, err.Error())
							continue
						}
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
