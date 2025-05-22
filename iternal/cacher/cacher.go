package cacher

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/vadim8q258475/store-product-microservice/config"
	repo "github.com/vadim8q258475/store-product-microservice/iternal/repo/sqlx/product"
)

type Cacher interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte) error
	Delete(ctx context.Context, key ...string) error
	DeleteAllByPrefix(ctx context.Context, prefix string) error
	DeleteListKeysByProductId(ctx context.Context, listPrefix string, id uint32) error
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

func (c *cacher) DeleteAllByPrefix(ctx context.Context, prefix string) error {
	var cursor uint64
	var keys []string
	var err error
	for {
		keys, cursor, err = c.client.Scan(ctx, cursor, prefix+"*", 100).Result()
		if err != nil {
			return err
		}

		if len(keys) > 0 {
			if err := c.client.Del(ctx, keys...).Err(); err != nil {
				return err
			}
		}

		if cursor == 0 {
			break
		}
	}

	return nil
}

func (c *cacher) DeleteListKeysByProductId(ctx context.Context, listPrefix string, id uint32) error {
	result := c.client.Keys(ctx, listPrefix+"*")
	if result.Err() != nil {
		return result.Err()
	}
	keys, err := result.Result()
	if err != nil {
		return err
	}
	wg := sync.WaitGroup{}
	ch := make(chan error)
	for _, key := range keys {
		wg.Add(1)
		go func(key string) {
			defer wg.Done()
			data, err := c.Get(ctx, key)
			if err != nil {
				select {
				case ch <- err:
				default:
				}
				return
			}
			var products []repo.Product
			err = json.Unmarshal(data, &products)
			if err != nil {
				select {
				case ch <- err:
				default:
				}
				return
			}
			for _, product := range products {
				if product.ID == id {
					err = c.Delete(ctx, key)
					if err != nil {
						select {
						case ch <- err:
						default:
						}
						return
					}
				}
			}
		}(key)
	}
	wg.Wait()
	select {
	case err = <-ch:
		return err
	default:
		return nil
	}
}
