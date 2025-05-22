package proxy

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/vadim8q258475/store-product-microservice/iternal/cacher"
	repo "github.com/vadim8q258475/store-product-microservice/iternal/repo/sqlx/product"
)

const listKeyPrefix = "list:"
const idKeyPrefix = "product:"

type productProxy struct {
	cacher cacher.Cacher
	repo   repo.ProductRepo
}

func NewProductProxy(cacher cacher.Cacher, repo repo.ProductRepo) repo.ProductRepo {
	return &productProxy{
		cacher: cacher,
		repo:   repo,
	}
}

func GetListKey(args []interface{}) (string, error) {
	data, err := json.Marshal(args)
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256(data)
	return listKeyPrefix + hex.EncodeToString(hash[:]), nil
}

func (p *productProxy) List(ctx context.Context, query string, args []interface{}) ([]repo.Product, error) {
	listKey, err := GetListKey(args)
	if err != nil {
		return nil, err
	}
	value, err := p.cacher.Get(ctx, listKey)
	if err != nil {
		if err == redis.Nil {
			result, err := p.repo.List(ctx, query, args)
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
	var result []repo.Product
	err = json.Unmarshal(value, &result)
	return result, err
}

func (p *productProxy) Create(ctx context.Context, product repo.Product) (uint32, error) {
	result, err := p.repo.Create(ctx, product)
	if err != nil {
		return result, err
	}
	err = p.cacher.DeleteListKeysByProductId(ctx, listKeyPrefix, result)
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
	err = p.cacher.DeleteListKeysByProductId(ctx, listKeyPrefix, id)
	if err != nil {
		return err
	}
	idKey := fmt.Sprintf("%s%d", idKeyPrefix, id)
	return p.cacher.Delete(ctx, idKey)
}
func (p *productProxy) Update(ctx context.Context, product repo.Product) (uint32, error) {
	result, err := p.repo.Update(ctx, product)
	if err != nil {
		return result, err
	}
	err = p.cacher.DeleteListKeysByProductId(ctx, listKeyPrefix, product.ID)
	if err != nil {
		return result, err
	}
	idKey := fmt.Sprintf("%s%d", idKeyPrefix, result)
	err = p.cacher.Delete(ctx, idKey)
	return result, err
}
func (p *productProxy) GetById(ctx context.Context, id uint32) (repo.Product, error) {
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
		return repo.Product{}, err
	}
	var result repo.Product
	err = json.Unmarshal(value, &result)
	return result, err
}
