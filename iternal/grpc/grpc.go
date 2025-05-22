package grpc

import (
	"context"

	gen "github.com/vadim8q258475/store-product-microservice/gen/v1"
	categoryService "github.com/vadim8q258475/store-product-microservice/iternal/service/category"
	productService "github.com/vadim8q258475/store-product-microservice/iternal/service/product"
)

type GrpcService struct {
	gen.UnimplementedProductServiceServer
	productService  productService.ProductService
	categoryService categoryService.CategoryService
}

func NewGrpcService(productService productService.ProductService, categoryService categoryService.CategoryService) *GrpcService {
	return &GrpcService{
		productService:  productService,
		categoryService: categoryService,
	}
}

func (g *GrpcService) CategoryCreate(ctx context.Context, request *gen.CategoryCreate_Request) (*gen.CategoryCreate_Response, error) {
	return g.categoryService.Create(ctx, request)
}
func (g *GrpcService) CategoryDelete(ctx context.Context, request *gen.CategoryDelete_Request) (*gen.CategoryDelete_Response, error) {
	return g.categoryService.Delete(ctx, request)
}
func (g *GrpcService) CategoryGetById(ctx context.Context, request *gen.CategoryGetById_Request) (*gen.CategoryGetById_Response, error) {
	return g.categoryService.GetById(ctx, request)
}
func (g *GrpcService) CategoryList(ctx context.Context, request *gen.CategoryList_Request) (*gen.CategoryList_Response, error) {
	return g.categoryService.List(ctx)
}
func (g *GrpcService) CategoryUpdate(ctx context.Context, request *gen.CategoryUpdate_Request) (*gen.CategoryUpdate_Response, error) {
	return g.categoryService.Update(ctx, request)
}
func (g *GrpcService) Create(ctx context.Context, request *gen.Create_Request) (*gen.Create_Response, error) {
	return g.productService.Create(ctx, request)
}
func (g *GrpcService) Delete(ctx context.Context, request *gen.Delete_Request) (*gen.Delete_Response, error) {
	return g.productService.Delete(ctx, request)
}
func (g *GrpcService) GetById(ctx context.Context, request *gen.GetById_Request) (*gen.GetById_Response, error) {
	return g.productService.GetById(ctx, request)
}
func (g *GrpcService) List(ctx context.Context, request *gen.List_Request) (*gen.List_Response, error) {
	return g.productService.List(ctx, request)
}
func (g *GrpcService) Update(ctx context.Context, request *gen.Update_Request) (*gen.Update_Response, error) {
	return g.productService.Update(ctx, request)
}
