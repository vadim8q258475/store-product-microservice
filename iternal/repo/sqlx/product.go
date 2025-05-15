package repo

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type Product struct {
	ID          uint32  `db:"id"`
	Name        string  `db:"name"`
	Description string  `db:"description"`
	Qty         int32   `db:"qty"`
	Price       float32 `db:"price"`
	CategoryID  uint32  `db:"category_id"`
}

type ProductRepo interface {
	List(ctx context.Context) ([]Product, error)
	Create(ctx context.Context, product Product) (uint32, error)
	Delete(ctx context.Context, id uint32) error
	Update(ctx context.Context, product Product) (uint32, error)
	GetById(ctx context.Context, id uint32) (Product, error)
}

type productRepo struct {
	db *sqlx.DB
}

func NewProductRepo(db *sqlx.DB) ProductRepo {
	return &productRepo{
		db: db,
	}
}

func (r *productRepo) List(ctx context.Context) ([]Product, error) {
	var products []Product
	err := r.db.SelectContext(ctx, &products, "SELECT id, name, description, qty, price, category_id FROM products")
	return products, err
}

func (r *productRepo) Create(ctx context.Context, product Product) (uint32, error) {
	var id uint32
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO products (name, description, qty, price, category_id)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		product.Name, product.Description, product.Qty, product.Price, product.CategoryID,
	).Scan(&id)
	return id, err
}

func (r *productRepo) Delete(ctx context.Context, id uint32) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM products WHERE id = $1", id)
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

func (r *productRepo) Update(ctx context.Context, product Product) (uint32, error) {
	var id uint32
	err := r.db.QueryRowContext(ctx,
		`UPDATE products
		SET name = $1, description = $2, qty = $3, price = $4, category_id = $5
		WHERE id = $6 RETURNING id`,
		product.Name, product.Description, product.Qty, product.Price, product.CategoryID, product.ID,
	).Scan(&id)
	return id, err
}

func (r *productRepo) GetById(ctx context.Context, id uint32) (Product, error) {
	var product Product
	err := r.db.GetContext(ctx, &product,
		"SELECT id, name, description, qty, price, category_id FROM products WHERE id = $1", id)
	return product, err
}
