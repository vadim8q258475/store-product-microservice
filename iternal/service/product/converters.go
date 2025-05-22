package service

import (
	"github.com/Masterminds/squirrel"
	gen "github.com/vadim8q258475/store-product-microservice/gen/v1"
	categoryRepo "github.com/vadim8q258475/store-product-microservice/iternal/repo/sqlx/category"
	productRepo "github.com/vadim8q258475/store-product-microservice/iternal/repo/sqlx/product"
)

func ProductCreateRequestToModel(request *gen.Create_Request) productRepo.Product {
	return productRepo.Product{
		Name:        request.Name,
		Description: request.Description,
		Qty:         request.Qty,
		Price:       request.Price,
		CategoryID:  request.CategoryId,
	}
}

func ProductUpdateRequestToModel(request *gen.Update_Request) productRepo.Product {
	return productRepo.Product{
		ID:          request.Id,
		Name:        request.Name,
		Description: request.Description,
		Qty:         request.Qty,
		Price:       request.Price,
		CategoryID:  request.CategoryId,
	}
}

func ProductModelToRequest(productModel productRepo.Product, categoryModel categoryRepo.Category) *gen.Product {
	return &gen.Product{
		Id:          productModel.ID,
		Name:        productModel.Name,
		Description: productModel.Description,
		Qty:         productModel.Qty,
		Price:       productModel.Price,
		Category: &gen.Category{
			Id:          categoryModel.ID,
			Name:        categoryModel.Name,
			Description: categoryModel.Description,
		},
	}
}

func MakeQuery(request *gen.List_Request) (string, []interface{}, error) {
	query := squirrel.Select("id", "name", "description", "qty", "price", "category_id").From("products")
	if len(request.CategoryIds) > 0 {
		query = query.Where(squirrel.Eq{"category_id": request.CategoryIds})
	}
	if request.MinPrice != 0 {
		query = query.Where(squirrel.Gt{"price": request.MinPrice})
	}
	if request.MaxPrice != 0 && request.MaxPrice > request.MinPrice {
		query = query.Where(squirrel.Lt{"price": request.MaxPrice})
	}
	if len(request.KeyWords) != 0 {
		orClauses := make([]squirrel.Sqlizer, 0)
		for _, key := range request.KeyWords {
			orClauses = append(orClauses, squirrel.Or{
				squirrel.Like{"name": "%" + key + "%"},
				squirrel.Like{"description": "%" + key + "%"},
			})
		}
		query = query.Where(squirrel.And(orClauses))
	}
	if request.SortBy == "name" || request.SortBy == "price" {
		orderBy := request.SortBy
		if request.Asc {
			orderBy += " ASC"
		} else {
			orderBy += " DESC"
		}
		query = query.OrderBy(orderBy)
	}
	offset := (request.Page - 1) * request.PageSize
	query.Limit(uint64(request.PageSize))
	query.Offset(uint64(offset))

	sql, args, err := query.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return "", []interface{}{}, err
	}
	return sql, args, nil
}
