package cmd

import (
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/api/notifier"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/api/trino"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/common/logging"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/configuration"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/discovery"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/proxy/healthcheck"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/proxy/session"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"time"
)

var rootCmd = &cobra.Command{
	Use:   "trino-loadbalancer",
	Short: "trino-loadbalancer is a fast, high available loadbalancer for trino",
}

var (
	configPath         string
	logger             logging.Logger = logging.Logrus()
	discoveryStorage   discovery.Storage
	sessionStorage     session.Storage
	clusterStats       trino.Api
	clusterHealthCheck healthcheck.HealthCheck
	discover           discovery.Discovery
	redisClient        redis.UniversalClient
	notifiers          notifier.Notifier
)

func init() {
	rootCmd.PersistentFlags().StringVar(&configPath, "config", "", "config file path")

	viper.SetDefault("persistence.postgres.db", "postgres")
	viper.SetDefault("persistence.postgres.host", "127.0.0.1")
	viper.SetDefault("persistence.postgres.port", 5432)
	viper.SetDefault("persistence.postgres.username", "postgres")
	viper.SetDefault("persistence.postgres.password", "")
	viper.SetDefault("persistence.postgres.ssl_mode", "disable")

	viper.SetDefault("session.store.redis.opts.prefix", "github.com/The-Data-Appeal-Company/trino-loadbalancer::")
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

	viper.SetDefault("controller.features.slow_worker_drainer.analyzer.std_deviation_ratio", 1.1)
	viper.SetDefault("controller.features.slow_worker_drainer.gracePeriodSeconds", 300)
	viper.SetDefault("controller.features.slow_worker_drainer.dryRun", true)
	viper.SetDefault("controller.features.slow_worker_drainer.drainThreshold", 3)

	cobra.OnInitialize(func() {
		err := readConfig()
		if err != nil {
			log.Fatal(err)
		}

		discoveryStorage, err = configuration.CreateDiscoveryStorage(configuration.DiscoveryStorageConfiguration{
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

		redisConfig := configuration.SessionStorageConfiguration{
			Standalone: configuration.RedisSessionStorageConfiguration{
				Enabled:  viper.GetBool("session.store.redis.standalone.enabled"),
				Host:     viper.GetString("session.store.redis.standalone.host"),
				DB:       viper.GetInt("session.store.redis.standalone.db"),
				Password: viper.GetString("session.store.redis.standalone.password"),
			},
			Sentinel: configuration.RedisSentinelSessionStorageConfiguration{
				Enabled:    viper.GetBool("session.store.redis.sentinel.enabled"),
				DB:         viper.GetInt("session.store.redis.sentinel.db"),
				Password:   viper.GetString("session.store.redis.sentinel.password"),
				MasterName: viper.GetString("session.store.redis.sentinel.master"),
				Hosts:      viper.GetStringSlice("session.store.redis.sentinel.hosts"),
			},
			Opts: configuration.RedisSessionStorageOpts{
				Prefix: viper.GetString("session.store.redis.opts.prefix"),
				MaxTTL: viper.GetDuration("session.store.redis.opts.max_ttl"),
			},
		}

		redisClient, err = configuration.CreateRedisStorageClient(redisConfig)
		if err != nil {
			log.Fatal(err)
		}

		sessionStorage, err = configuration.CreateSessionStorage(redisClient, redisConfig)

		if err != nil {
			log.Fatal(err)
		}

		clusterHealthCheck, err = configuration.CreateHealthCheck(configuration.HealthCheckConfiguration{
			Enabled: viper.GetBool("clusters.healthcheck.enabled"),
		})
		if err != nil {
			log.Fatal(err)
		}

		clusterStats, err = configuration.CreateStatisticsRetriever(configuration.StatisticsConfiguration{
			Enabled: viper.GetBool("clusters.statistics.enabled"),
		})

		if err != nil {
			log.Fatal(err)
		}

		var conf []configuration.DiscoveryConfiguration
		err = viper.UnmarshalKey("discovery.providers", &conf)
		if err != nil {
			log.Fatal(err)
		}

		discover, err = configuration.CreateCrossProviderDiscovery(conf)
		if err != nil {
			log.Fatal(err)
		}

		var notifierConfig configuration.NotifierConfig
		err = viper.UnmarshalKey("notifier", &notifierConfig)
		if err != nil {
			log.Fatal(err)
		}

		notifiers = configuration.CreateNotifier(notifierConfig)

	})
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
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
