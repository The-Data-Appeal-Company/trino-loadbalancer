package session

import (
	"context"
	"github.com/The-Data-Appeal-Company/presto-loadbalancer/pkg/models"
)

func NewStorageCache(source Storage, cache *Memory) *StorageCache {
	return &StorageCache{
		cache:  cache,
		source: source,
	}
}

type StorageCache struct {
	cache  *Memory
	source Storage
}

func NewCaching(source Storage, cache *Memory) *StorageCache {
	return &StorageCache{
		source: source,
		cache:  cache,
	}
}

func (s StorageCache) Link(ctx context.Context, info models.QueryInfo, coordinator string) error {
	if err := s.cache.Link(ctx, info, coordinator); err != nil {
		return err
	}
	if err := s.source.Link(ctx, info, coordinator); err != nil {
		_ = s.cache.Unlink(ctx, info) // in case of error we rollback the cache to avoid state differences
		return err
	}
	return nil
}

func (s StorageCache) Unlink(ctx context.Context, info models.QueryInfo) error {
	if err := s.cache.Unlink(ctx, info); err != nil {
		return err
	}

	if err := s.source.Unlink(ctx, info); err != nil {
		// here we don't a fast way to rollback the cache but keeping a link to a terminated query
		// is not a problem so we prefer faster unlink
		return err
	}

	return nil
}

func (s StorageCache) Get(ctx context.Context, info models.QueryInfo) (string, error) {
	cached, err := s.cache.Get(ctx, info)
	if err != nil {
		if err != ErrLinkNotFound {
			return "", err
		}
	} else {
		return cached, nil
	}

	value, err := s.source.Get(ctx, info)
	if err != nil {
		return "", err
	}

	if err := s.cache.Link(ctx, info, value); err != nil {
		return "", err // we may want to ignore this
	}

	return value, nil
}
