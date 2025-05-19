package repo

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/vadim8q258475/store-product-microservice/iternal/cacher"
)

const listKey = "list"
const idKeyPrefix = "category:"

type categoryProxy struct {
	cacher cacher.Cacher
	repo   CategoryRepo
}

func NewCategoryProxy(cacher cacher.Cacher, repo CategoryRepo) CategoryRepo {
	return &categoryProxy{
		cacher: cacher,
		repo:   repo,
	}
}

func (p *categoryProxy) List(ctx context.Context) ([]Category, error) {
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
	var result []Category
	err = json.Unmarshal(value, &result)
	return result, err
}
func (p *categoryProxy) Create(ctx context.Context, category Category) (uint32, error) {
	result, err := p.repo.Create(ctx, category)
	if err != nil {
		return result, err
	}
	err = p.cacher.Delete(ctx, listKey)
	if err != nil {
		return result, err
	}
	return result, err
}
func (p *categoryProxy) Delete(ctx context.Context, id uint32) error {
	err := p.repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	idKey := fmt.Sprintf("%s%d", idKeyPrefix, id)
	return p.cacher.Delete(ctx, listKey, idKey)
}
func (p *categoryProxy) Update(ctx context.Context, category Category) (uint32, error) {
	result, err := p.repo.Update(ctx, category)
	if err != nil {
		return result, err
	}
	idKey := fmt.Sprintf("%s%d", idKeyPrefix, result)
	err = p.cacher.Delete(ctx, listKey, idKey)
	if err != nil {
		return result, err
	}
	return result, err
}
func (p *categoryProxy) GetById(ctx context.Context, id uint32) (Category, error) {
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
		return Category{}, err
	}
	var result Category
	err = json.Unmarshal(value, &result)
	return result, err
}
