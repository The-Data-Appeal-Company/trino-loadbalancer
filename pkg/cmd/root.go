package cmd

import (
	"fmt"
	"github.com/The-Data-Appeal-Company/presto-loadbalancer/pkg/discovery"
	"github.com/The-Data-Appeal-Company/presto-loadbalancer/pkg/factory"
	"github.com/The-Data-Appeal-Company/presto-loadbalancer/pkg/healthcheck"
	"github.com/The-Data-Appeal-Company/presto-loadbalancer/pkg/logging"
	"github.com/The-Data-Appeal-Company/presto-loadbalancer/pkg/session"
	"github.com/The-Data-Appeal-Company/presto-loadbalancer/pkg/statistics"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"time"
)

var rootCmd = &cobra.Command{
	Use:   "github.com/The-Data-Appeal-Company/presto-loadbalancer",
	Short: "github.com/The-Data-Appeal-Company/presto-loadbalancer is a fast, high available loadbalancer for presto",
}

var (
	configPath         string
	logger             logging.Logger = logging.Logrus()
	discoveryStorage   discovery.Storage
	sessionStorage     session.Storage
	clusterStats       statistics.Retriever
	clusterHealthCheck healthcheck.HealthCheck
	discover           discovery.Discovery
)

type DiscoveryConf struct {
	provider string            `json:"provider"`
	enabled  bool              `json:"enabled"`
	aws      AwsProviderParams `json:"aws"`
}

type AwsProviderParams struct {
	accessKeyId string `json:"access_key_id"`
	secretKey   string `json:"secret_key"`
	region      string `json:"region"`
}

func init() {
	rootCmd.PersistentFlags().StringVar(&configPath, "config", "", "config file path")

	viper.SetDefault("persistence.postgres.db", "postgres")
	viper.SetDefault("persistence.postgres.host", "127.0.0.1")
	viper.SetDefault("persistence.postgres.port", 5432)
	viper.SetDefault("persistence.postgres.username", "postgres")
	viper.SetDefault("persistence.postgres.password", "")
	viper.SetDefault("persistence.postgres.ssl_mode", "disable")

	viper.SetDefault("session.store.redis.opts.prefix", "github.com/The-Data-Appeal-Company/presto-loadbalancer::")
	viper.SetDefault("session.store.redis.opts.max_ttl", 24*time.Hour)

	viper.SetDefault("session.store.redis.standalone.enabled", true)
	viper.SetDefault("session.store.redis.standalone.host", "127.0.0.1:6379")
	viper.SetDefault("session.store.redis.standalone.db", 0)

	viper.SetDefault("session.store.redis.sentinel.enabled", false)
	viper.SetDefault("session.store.redis.sentinel.hosts", []string{"127.0.0.1:26379"})
	viper.SetDefault("session.store.redis.sentinel.db", 0)
	viper.SetDefault("session.store.redis.sentinel.master", "mymaster")

	viper.SetDefault("clusters.statistics.enabled", true)
	viper.SetDefault("clusters.healthcheck.enabled", true)

	viper.SetDefault("proxy.port", 8998)

	viper.SetDefault("routing.rule", "round-robin")

	viper.SetDefault("clusters.healthcheck.delay", 10*time.Second)
	viper.SetDefault("clusters.statistics.delay", 10*time.Second)
	viper.SetDefault("clusters.sync.delay", 10*time.Minute)

	viper.SetDefault("discovery.enabled", false)

	cobra.OnInitialize(func() {
		err := readConfig()
		if err != nil {
			log.Fatal(err)
		}

		discoveryStorage, err = factory.CreateDiscoveryStorage(factory.DiscoveryStorageConfiguration{
			Db:       viper.GetString("persistence.postgres.db"),
			Host:     viper.GetString("persistence.postgres.host"),
			Port:     viper.GetInt("persistence.postgres.port"),
			User:     viper.GetString("persistence.postgres.username"),
			Password: viper.GetString("persistence.postgres.password"),
			SslMode:  viper.GetString("persistence.postgres.ssl_mode"),
		})

		if err != nil {
			log.Fatal(err)
		}

		sessionStorage, err = factory.CreateSessionStorage(factory.SessionStorageConfiguration{
			Standalone: factory.RedisSessionStorageConfiguration{
				Enabled:  viper.GetBool("session.store.redis.standalone.enabled"),
				Host:     viper.GetString("session.store.redis.standalone.host"),
				DB:       viper.GetInt("session.store.redis.standalone.db"),
				Password: viper.GetString("session.store.redis.standalone.password"),
			},
			Sentinel: factory.RedisSentinelSessionStorageConfiguration{
				Enabled:    viper.GetBool("session.store.redis.sentinel.enabled"),
				DB:         viper.GetInt("session.store.redis.sentinel.db"),
				Password:   viper.GetString("session.store.redis.sentinel.password"),
				MasterName: viper.GetString("session.store.redis.sentinel.master"),
				Hosts:      viper.GetStringSlice("session.store.redis.sentinel.hosts"),
			},
			Opts: factory.RedisSessionStorageOpts{
				Prefix: viper.GetString("session.store.redis.opts.prefix"),
				MaxTTL: viper.GetDuration("session.store.redis.opts.max_ttl"),
			},
		})

		if err != nil {
			log.Fatal(err)
		}

		clusterHealthCheck, err = factory.CreateHealthCheck(factory.HealthCheckConfiguration{
			Enabled: viper.GetBool("clusters.healthcheck.enabled"),
		})
		if err != nil {
			log.Fatal(err)
		}

		clusterStats, err = factory.CreateStatisticsRetriever(factory.StatisticsConfiguration{
			Enabled: viper.GetBool("clusters.statistics.enabled"),
		})

		if err != nil {
			log.Fatal(err)
		}

		var discoveryConfs []factory.DiscoveryConfiguration
		err = viper.UnmarshalKey("discovery.providers", discoveryConfs)

		if err != nil {
			log.Fatal(err)
		}

		discover, err = factory.CreateCrossProviderDiscovery(discoveryConfs, rootCmd.Context())

		if err != nil {
			log.Fatal(err)
		}

	})
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func readConfig() error {
	if configPath != "" {
		viper.SetConfigFile(configPath)
	} else {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("/config")
		viper.AddConfigPath(".")
	}

	viper.AutomaticEnv()

	return viper.ReadInConfig()
}
