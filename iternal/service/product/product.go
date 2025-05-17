package service

import (
	"context"
	"sync"

	gen "github.com/vadim8q258475/store-product-microservice/gen/v1"
	repo "github.com/vadim8q258475/store-product-microservice/iternal/repo/sqlx"
)

type ProductService interface {
	List(ctx context.Context) (*gen.List_Response, error)
	Create(ctx context.Context, request *gen.Create_Request) (*gen.Create_Response, error)
	Delete(ctx context.Context, request *gen.Delete_Request) (*gen.Delete_Response, error)
	Update(ctx context.Context, request *gen.Update_Request) (*gen.Update_Response, error)
	GetById(ctx context.Context, request *gen.GetById_Request) (*gen.GetById_Response, error)
}

type productService struct {
	productRepo  repo.ProductRepo
	categoryRepo repo.CategoryRepo
}

func NewProductService(productRepo repo.ProductRepo, categoryRepo repo.CategoryRepo) ProductService {
	return &productService{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
	}
}

func (s *productService) List(ctx context.Context) (*gen.List_Response, error) {
	productModels, err := s.productRepo.List(ctx)
	if err != nil {
		return nil, err
	}
	products := make([]*gen.Product, len(productModels))
	wg := sync.WaitGroup{}
	errCh := make(chan error)
	for i, productModel := range productModels {
		wg.Add(1)
		go func(i int, productModel repo.Product) {
			defer wg.Done()
			categoryModel, err := s.categoryRepo.GetById(ctx, productModel.CategoryID)
			if err != nil {
				select {
				case errCh <- err:
				default:
				}
				return
			}
			products[i] = ProductModelToRequest(productModel, categoryModel)
		}(i, productModel)
	}
	wg.Wait()
	select {
	case err = <-errCh:
		return nil, err
	default:
		return &gen.List_Response{Products: products}, nil
	}
}

func (s *productService) Create(ctx context.Context, request *gen.Create_Request) (*gen.Create_Response, error) {
	model := ProductCreateRequestToModel(request)
	id, err := s.productRepo.Create(ctx, model)
	return &gen.Create_Response{Id: id}, err
}

func (s *productService) Delete(ctx context.Context, request *gen.Delete_Request) (*gen.Delete_Response, error) {
	err := s.productRepo.Delete(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	return &gen.Delete_Response{Success: true}, nil
}

func (s *productService) Update(ctx context.Context, request *gen.Update_Request) (*gen.Update_Response, error) {
	model := ProductUpdateRequestToModel(request)
	id, err := s.productRepo.Update(ctx, model)
	return &gen.Update_Response{Id: id}, err
}

func (s *productService) GetById(ctx context.Context, request *gen.GetById_Request) (*gen.GetById_Response, error) {
	productModel, err := s.productRepo.GetById(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	categoryModel, err := s.categoryRepo.GetById(ctx, productModel.CategoryID)
	return &gen.GetById_Response{
		Product: ProductModelToRequest(productModel, categoryModel),
	}, err
}
