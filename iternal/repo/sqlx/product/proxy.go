package repo

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/vadim8q258475/store-product-microservice/iternal/cacher"
)

const listKey = "list"
const idKeyPrefix = "product:"

type productProxy struct {
	cacher cacher.Cacher
	repo   ProductRepo
}

func NewProductProxy(cacher cacher.Cacher, repo ProductRepo) ProductRepo {
	return &productProxy{
		cacher: cacher,
		repo:   repo,
	}
}

func (p *productProxy) List(ctx context.Context) ([]Product, error) {
	value, err := p.cacher.Get(ctx, listKey)
	if err != nil {
		if err == redis.Nil {
			result, err := p.repo.List(ctx)
			if err != nil {
				return nil, err
			}
			bytes, err := json.Marshal(result)
			if err != nil {
				return nil, err
			}
			err = p.cacher.Set(ctx, listKey, bytes)
			if err != nil {
				return nil, err
			}
			return result, nil
		}
		return nil, err
	}
	var result []Product
	err = json.Unmarshal(value, &result)
	return result, err
}
func (p *productProxy) Create(ctx context.Context, product Product) (uint32, error) {
	result, err := p.repo.Create(ctx, product)
	if err != nil {
		return result, err
	}
	err = p.cacher.Delete(ctx, listKey)
	if err != nil {
		return result, err
	}
	return result, err
}
func (p *productProxy) Delete(ctx context.Context, id uint32) error {
	err := p.repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	idKey := fmt.Sprintf("%s%d", idKeyPrefix, id)
	return p.cacher.Delete(ctx, listKey, idKey)
}
func (p *productProxy) Update(ctx context.Context, product Product) (uint32, error) {
	result, err := p.repo.Update(ctx, product)
	if err != nil {
		return result, err
	}
	idKey := fmt.Sprintf("%s%d", idKeyPrefix, result)
	err = p.cacher.Delete(ctx, listKey, idKey)
	return result, err
}
func (p *productProxy) GetById(ctx context.Context, id uint32) (Product, error) {
	key := fmt.Sprintf("%s%d", idKeyPrefix, id)
	value, err := p.cacher.Get(ctx, key)
	if err != nil {
		if err == redis.Nil {
			result, err := p.repo.GetById(ctx, id)
			if err != nil {
				return result, err
			}
			bytes, err := json.Marshal(result)
			if err != nil {
				return result, err
			}
			err = p.cacher.Set(ctx, key, bytes)
			if err != nil {
				return result, err
			}
			return result, nil
		}
		return Product{}, err
	}
	var result Product
	err = json.Unmarshal(value, &result)
	return result, err
}
