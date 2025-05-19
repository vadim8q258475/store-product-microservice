package service

import (
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
