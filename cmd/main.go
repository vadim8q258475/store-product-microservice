package main

import (
	"fmt"

	"github.com/vadim8q258475/store-product-microservice/app"
	"github.com/vadim8q258475/store-product-microservice/config"
	grpcService "github.com/vadim8q258475/store-product-microservice/iternal/grpc"
	"github.com/vadim8q258475/store-product-microservice/iternal/interceptor"
	repo "github.com/vadim8q258475/store-product-microservice/iternal/repo/sqlx"

	categoryService "github.com/vadim8q258475/store-product-microservice/iternal/service/category"
	productService "github.com/vadim8q258475/store-product-microservice/iternal/service/product"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// TODO
// add cacher

func main() {
	// logger
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	// interceptor
	intterceptor := interceptor.NewInterceptor(logger)

	// load config
	cfg := config.MustLoadConfig()
	fmt.Println(cfg.String())

	// init db
	client, err := repo.InitDB(cfg)
	if err != nil {
		panic(err)
	}

	defer client.Close()

	// repo
	categoryRepo := repo.NewCategoryRepo(client)
	productRepo := repo.NewProductRepo(client)

	// service
	productService := productService.NewProductService(productRepo, categoryRepo)
	categoryService := categoryService.NewCategoryService(categoryRepo)

	// grpc service
	grpcService := grpcService.NewGrpcService(productService, categoryService)

	// grpc server
	server := grpc.NewServer(
		grpc.UnaryInterceptor(
			intterceptor.UnaryServerInterceptor,
		),
	)

	// app
	app := app.NewApp(grpcService, server, logger, cfg)

	if err = app.Run(); err != nil {
		panic(err)
	}
}
