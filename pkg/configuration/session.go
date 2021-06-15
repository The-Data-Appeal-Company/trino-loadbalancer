package configuration

import (
	"context"
	"errors"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/proxy/session"
	"github.com/go-redis/redis/v8"
	"time"
)

type RedisSessionStorageConfiguration struct {
	Enabled  bool
	Host     string
	DB       int
	Password string
}

type RedisSentinelSessionStorageConfiguration struct {
	Enabled    bool
	Hosts      []string
	MasterName string
	DB         int
	Password   string
}

type RedisSessionStorageOpts struct {
	Prefix string
	MaxTTL time.Duration
}

type SessionStorageConfiguration struct {
	Standalone RedisSessionStorageConfiguration
	Sentinel   RedisSentinelSessionStorageConfiguration
	Opts       RedisSessionStorageOpts
}

func CreateSessionStorage(conf SessionStorageConfiguration) (session.Storage, error) {
	if !conf.Standalone.Enabled && !conf.Sentinel.Enabled {
		return nil, errors.New("no session storage enabled")
	}

	var client redis.UniversalClient

	if conf.Standalone.Enabled {
		client = redis.NewClient(&redis.Options{
			Addr:     conf.Standalone.Host,
			Password: conf.Standalone.Password,
			DB:       conf.Standalone.DB,
		})
	}

	if conf.Sentinel.Enabled {
		client = redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:    conf.Sentinel.MasterName,
			SentinelAddrs: conf.Sentinel.Hosts,
			DB:            conf.Sentinel.DB,
			Password:      conf.Sentinel.Password,
		})
	}

	_, err := client.Ping(context.TODO()).Result()
	if err != nil {
		return nil, err
	}

	redisStorage := session.NewRedisStorage(client, conf.Opts.Prefix, conf.Opts.MaxTTL)
	memoryStorage := session.NewMemoryStorage()
	return session.NewStorageCache(redisStorage, memoryStorage), nil
}
