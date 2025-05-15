package repo

import (
	"context"

	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	ID          uint32 `gorm:"primarykey"`
	Name        string
	Description string
	Products    []Product
}

type CategoryRepo interface {
	List(ctx context.Context) ([]Category, error)
	Create(ctx context.Context, category Category) (uint32, error)
	Delete(ctx context.Context, id uint32) error
	Update(ctx context.Context, category Category) (uint32, error)
	GetById(ctx context.Context, id uint32) (Category, error)
}

type categoryRepo struct {
	db *gorm.DB
}

func NewCategoryRepo(db *gorm.DB) CategoryRepo {
	return &categoryRepo{
		db: db,
	}
}

func (r *categoryRepo) List(ctx context.Context) ([]Category, error) {
	var categories []Category
	result := r.db.WithContext(ctx).Find(&categories)
	return categories, result.Error
}

func (r *categoryRepo) Create(ctx context.Context, category Category) (uint32, error) {
	result := r.db.WithContext(ctx).Create(&category)
	return category.ID, result.Error
}

func (r *categoryRepo) GetById(ctx context.Context, id uint32) (Category, error) {
	category := Category{ID: id}
	result := r.db.WithContext(ctx).Find(&category)
	return category, result.Error
}

func (r *categoryRepo) Delete(ctx context.Context, id uint32) error {
	category := Category{ID: id}
	result := r.db.WithContext(ctx).Delete(&category)
	return result.Error
}

func (r *categoryRepo) Update(ctx context.Context, category Category) (uint32, error) {
	result := r.db.WithContext(ctx).Save(&category)
	return category.ID, result.Error
}
