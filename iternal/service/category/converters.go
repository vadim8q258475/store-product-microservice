package service

import (
	gen "github.com/vadim8q258475/store-product-microservice/gen/v1"
	repo "github.com/vadim8q258475/store-product-microservice/iternal/repo/sqlx/category"
)

func CategoryRequestToModel(request *gen.Category) repo.Category {
	return repo.Category{
		ID:          request.Id,
		Name:        request.Name,
		Description: request.Description,
	}
}

func CategoryCreateRequestToModel(request *gen.CategoryCreate_Request) repo.Category {
	return repo.Category{
		Name:        request.Name,
		Description: request.Description,
	}
}

func CategoryUpdateRequestToModel(request *gen.CategoryUpdate_Request) repo.Category {
	return repo.Category{
		ID:          request.Id,
		Name:        request.Name,
		Description: request.Description,
	}
}

func CategoryModelToRequest(model repo.Category) *gen.Category {
	return &gen.Category{
		Id:          model.ID,
		Name:        model.Name,
		Description: model.Description,
	}
}
