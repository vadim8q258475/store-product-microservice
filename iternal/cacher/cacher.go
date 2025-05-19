package cacher

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/vadim8q258475/store-product-microservice/config"
)

type Cacher interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte) error
	Delete(ctx context.Context, key ...string) error
}

type cacher struct {
	client       *redis.Client
	cacheMinutes int
}

func NewCacher(client *redis.Client, cfg config.Config) Cacher {
	return &cacher{
		client:       client,
		cacheMinutes: cfg.CacheMinutes,
	}
}

func (c *cacher) Get(ctx context.Context, key string) ([]byte, error) {
	value, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (c *cacher) Set(ctx context.Context, key string, value []byte) error {
	err := c.client.Set(ctx, key, value, time.Minute*time.Duration(c.cacheMinutes)).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *cacher) Delete(ctx context.Context, keys ...string) error {
	err := c.client.Del(ctx, keys...).Err()
	if err != nil {
		return err
	}
	return nil
}
