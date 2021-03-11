package session

import (
	"context"
	"fmt"
	"github.com/The-Data-Appeal-Company/trino-loadbalancer/pkg/models"
	"github.com/go-redis/redis/v8"
	"time"
)

func NewRedisStorage(client redis.UniversalClient, prefix string, ttl time.Duration) RedisLinkerStorage {
	return RedisLinkerStorage{
		redis:  client,
		prefix: prefix,
		ttl:    ttl,
	}
}

type RedisLinkerStorage struct {
	redis  redis.UniversalClient
	prefix string
	ttl    time.Duration
}

func (r RedisLinkerStorage) Link(ctx context.Context, info models.QueryInfo, coordinator string) error {
	return r.redis.Set(ctx, r.queryHash(info), coordinator, r.ttl).Err()
}

func (r RedisLinkerStorage) Unlink(ctx context.Context, info models.QueryInfo) error {
	return r.redis.Del(ctx, r.queryHash(info)).Err()
}

func (r RedisLinkerStorage) Get(ctx context.Context, info models.QueryInfo) (string, error) {
	coordinator, err := r.redis.Get(ctx, r.queryHash(info)).Result()
	if err == redis.Nil {
		return "", ErrLinkNotFound
	}

	if err != nil {
		return "", err
	}
	return coordinator, nil
}

func (r RedisLinkerStorage) queryHash(info models.QueryInfo) string {
	return fmt.Sprintf("%s::%s::%s", r.prefix, info.TransactionID, info.QueryID)
}
