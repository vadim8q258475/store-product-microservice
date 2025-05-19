package repo

import (
	"context"
	"errors"

	"github.com/jmoiron/sqlx"
)

var ErrNotFound = errors.New("category not found")

type Category struct {
	ID          uint32 `db:"id"`
	Name        string `db:"name"`
	Description string `db:"description"`
}

type CategoryRepo interface {
	List(ctx context.Context) ([]Category, error)
	Create(ctx context.Context, category Category) (uint32, error)
	Delete(ctx context.Context, id uint32) error
	Update(ctx context.Context, category Category) (uint32, error)
	GetById(ctx context.Context, id uint32) (Category, error)
}

type categoryRepo struct {
	db *sqlx.DB
}

func NewCategoryRepo(db *sqlx.DB) CategoryRepo {
	return &categoryRepo{
		db: db,
	}
}

func (r *categoryRepo) List(ctx context.Context) ([]Category, error) {
	var categories []Category
	err := r.db.SelectContext(ctx, &categories, "SELECT id, name, description FROM categories")
	return categories, err
}

func (r *categoryRepo) Create(ctx context.Context, category Category) (uint32, error) {
	var id uint32
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO categories (name, description) 
		VALUES ($1, $2) RETURNING id`,
		category.Name, category.Description).Scan(&id)
	return id, err
}

func (r *categoryRepo) GetById(ctx context.Context, id uint32) (Category, error) {
	var category Category
	err := r.db.GetContext(ctx, &category,
		"SELECT id, name, description FROM categories WHERE id = $1", id)
	return category, err
}

func (r *categoryRepo) Delete(ctx context.Context, id uint32) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM categories WHERE id = $1", id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return err
}

func (r *categoryRepo) Update(ctx context.Context, category Category) (uint32, error) {
	var updatedID uint32
	err := r.db.QueryRowContext(ctx,
		`UPDATE categories 
         SET name = $1, description = $2 
         WHERE id = $3
         RETURNING id`,
		category.Name, category.Description, category.ID).Scan(&updatedID)
	return updatedID, err
}
