package service

import (
	gen "github.com/vadim8q258475/store-product-microservice/gen/v1"
	repo "github.com/vadim8q258475/store-product-microservice/iternal/repo/sqlx"
)

func ProductCreateRequestToModel(request *gen.Create_Request) repo.Product {
	return repo.Product{
		Name:        request.Name,
		Description: request.Description,
		Qty:         request.Qty,
		Price:       request.Price,
		CategoryID:  request.CategoryId,
	}
}

func ProductUpdateRequestToModel(request *gen.Update_Request) repo.Product {
	return repo.Product{
		ID:          request.Id,
		Name:        request.Name,
		Description: request.Description,
		Qty:         request.Qty,
		Price:       request.Price,
		CategoryID:  request.CategoryId,
	}
}

func ProductModelToRequest(productModel repo.Product, categoryModel repo.Category) *gen.Product {
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
