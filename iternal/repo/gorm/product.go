package repo

import (
	"context"

	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	ID          uint32 `gorm:"primarykey"`
	Name        string
	Description string
	Qty         int32
	Price       float64
	CategoryID  uint32
	Category    Category
}

type ProductRepo interface {
	List(ctx context.Context) ([]Product, error)
	Create(ctx context.Context, product Product) (uint32, error)
	Delete(ctx context.Context, id uint32) error
	Update(ctx context.Context, product Product) (uint32, error)
	GetById(ctx context.Context, id uint32) (Product, error)
}

type productRepo struct {
	db *gorm.DB
}

func NewproductRepo(db *gorm.DB) ProductRepo {
	return &productRepo{
		db: db,
	}
}

func (r *productRepo) List(ctx context.Context) ([]Product, error) {
	var categories []Product
	result := r.db.WithContext(ctx).Find(&categories)
	return categories, result.Error
}

func (r *productRepo) Create(ctx context.Context, product Product) (uint32, error) {
	result := r.db.WithContext(ctx).Create(&product)
	return product.ID, result.Error
}

func (r *productRepo) GetById(ctx context.Context, id uint32) (Product, error) {
	product := Product{ID: id}
	result := r.db.WithContext(ctx).Find(&product)
	return product, result.Error
}

func (r *productRepo) Delete(ctx context.Context, id uint32) error {
	product := Product{ID: id}
	result := r.db.WithContext(ctx).Delete(&product)
	return result.Error
}

func (r *productRepo) Update(ctx context.Context, product Product) (uint32, error) {
	result := r.db.WithContext(ctx).Save(&product)
	return product.ID, result.Error
}
