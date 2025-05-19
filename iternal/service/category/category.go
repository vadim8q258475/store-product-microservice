package service

import (
	"context"

	repo "github.com/vadim8q258475/store-product-microservice/iternal/repo/sqlx/category"

	gen "github.com/vadim8q258475/store-product-microservice/gen/v1"
)

type CategoryService interface {
	List(ctx context.Context) (*gen.CategoryList_Response, error)
	Create(ctx context.Context, request *gen.CategoryCreate_Request) (*gen.CategoryCreate_Response, error)
	Delete(ctx context.Context, request *gen.CategoryDelete_Request) (*gen.CategoryDelete_Response, error)
	Update(ctx context.Context, request *gen.CategoryUpdate_Request) (*gen.CategoryUpdate_Response, error)
	GetById(ctx context.Context, request *gen.CategoryGetById_Request) (*gen.CategoryGetById_Response, error)
}

type categoryService struct {
	repo repo.CategoryRepo
}

func NewCategoryService(repo repo.CategoryRepo) CategoryService {
	return &categoryService{
		repo: repo,
	}
}

func (s *categoryService) List(ctx context.Context) (*gen.CategoryList_Response, error) {
	models, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	categories := make([]*gen.Category, len(models))
	for i, model := range models {
		categories[i] = CategoryModelToRequest(model)
	}
	return &gen.CategoryList_Response{Categories: categories}, err
}

func (s *categoryService) Create(ctx context.Context, request *gen.CategoryCreate_Request) (*gen.CategoryCreate_Response, error) {
	model := CategoryCreateRequestToModel(request)
	id, err := s.repo.Create(ctx, model)
	return &gen.CategoryCreate_Response{Id: id}, err
}

func (s *categoryService) Delete(ctx context.Context, request *gen.CategoryDelete_Request) (*gen.CategoryDelete_Response, error) {
	err := s.repo.Delete(ctx, request.Id)
	if err != nil {
		return &gen.CategoryDelete_Response{Success: false}, err
	}
	return &gen.CategoryDelete_Response{Success: true}, nil
}

func (s *categoryService) Update(ctx context.Context, request *gen.CategoryUpdate_Request) (*gen.CategoryUpdate_Response, error) {
	model := CategoryUpdateRequestToModel(request)
	id, err := s.repo.Update(ctx, model)
	return &gen.CategoryUpdate_Response{Id: id}, err
}

func (s *categoryService) GetById(ctx context.Context, request *gen.CategoryGetById_Request) (*gen.CategoryGetById_Response, error) {
	model, err := s.repo.GetById(ctx, request.Id)
	return &gen.CategoryGetById_Response{Category: CategoryModelToRequest(model)}, err
}
