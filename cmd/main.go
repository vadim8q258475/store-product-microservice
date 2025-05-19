package main

import (
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/vadim8q258475/store-product-microservice/app"
	"github.com/vadim8q258475/store-product-microservice/config"
	"github.com/vadim8q258475/store-product-microservice/iternal/cacher"
	grpcService "github.com/vadim8q258475/store-product-microservice/iternal/grpc"
	"github.com/vadim8q258475/store-product-microservice/iternal/interceptor"
	repo "github.com/vadim8q258475/store-product-microservice/iternal/repo/sqlx"

	categoryRepo "github.com/vadim8q258475/store-product-microservice/iternal/repo/sqlx/category"
	productRepo "github.com/vadim8q258475/store-product-microservice/iternal/repo/sqlx/product"

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
	categoryRepos := categoryRepo.NewCategoryRepo(client)
	productRepos := productRepo.NewProductRepo(client)

	// init cache
	cacheClient := redis.NewClient(&redis.Options{
		Addr: cfg.CacheHost + ":" + cfg.CachePort,
	})
	cacher := cacher.NewCacher(cacheClient, cfg)

	// repo proxy
	categoryProxy := categoryRepo.NewCategoryProxy(cacher, categoryRepos)
	productProxy := productRepo.NewProductProxy(cacher, productRepos)

	// service
	productService := productService.NewProductService(productProxy, categoryProxy)
	categoryService := categoryService.NewCategoryService(categoryProxy)

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
