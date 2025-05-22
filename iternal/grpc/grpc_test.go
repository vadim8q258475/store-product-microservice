package grpc

import (
	"testing"
	// gen "github.com/vadim8q258475/store-product-microservice/gen/v1"
	// categoryRepo "github.com/vadim8q258475/store-product-microservice/iternal/repo/sqlx/category"
	// productRepo "github.com/vadim8q258475/store-product-microservice/iternal/repo/sqlx/product"
	// categoryService "github.com/vadim8q258475/store-product-microservice/iternal/service/category"
	// productService "github.com/vadim8q258475/store-product-microservice/iternal/service/product"
	// "go.uber.org/mock/gomock"
)

func TestGrpcService_List(t *testing.T) {

}

func TestGrpcService_Create(t *testing.T) {
	// ctrl := gomock.NewController(t)
	// defer ctrl.Finish()
	// productMockRepo := productRepo.NewMockProductRepo(ctrl)
	// categoryMockRepo := categoryRepo.NewMockCategoryRepo(ctrl)
	// categoryService := categoryService.NewCategoryService(categoryMockRepo)
	// productService := productService.NewProductService(productMockRepo, categoryMockRepo)
	// grpcService := NewGrpcService(productService, categoryService)

	// ctx := context.Background()
	// id := uint32(1)
	// productRepoModel := productRepo.Product{}
	// productCreateRequest := &gen.Create_Request{}
	// // all ok

	// productMockRepo.EXPECT().Create(ctx, productRepoModel).Return(id, nil)
}

func TestGrpcService_Delete(t *testing.T) {

}

func TestGrpcService_Update(t *testing.T) {

}

func TestGrpcService_GetByID(t *testing.T) {

}
