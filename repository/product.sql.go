// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: product.sql

package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createProduct = `-- name: CreateProduct :one
INSERT INTO products (
  store_id, sku, name, description, price, stock_quantity, created_at, updated_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING id, store_id, sku, name, description, price, stock_quantity, created_at, updated_at
`

type CreateProductParams struct {
	StoreID       uuid.UUID `json:"storeId"`
	Sku           string    `json:"sku"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	Price         string    `json:"price"`
	StockQuantity int32     `json:"stockQuantity"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

// Create a new product
func (q *Queries) CreateProduct(ctx context.Context, arg CreateProductParams) (Product, error) {
	row := q.db.QueryRowContext(ctx, createProduct,
		arg.StoreID,
		arg.Sku,
		arg.Name,
		arg.Description,
		arg.Price,
		arg.StockQuantity,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i Product
	err := row.Scan(
		&i.ID,
		&i.StoreID,
		&i.Sku,
		&i.Name,
		&i.Description,
		&i.Price,
		&i.StockQuantity,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteProduct = `-- name: DeleteProduct :exec
DELETE FROM products
WHERE id = $1
`

// Delete a product by ID
func (q *Queries) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteProduct, id)
	return err
}

const getProduct = `-- name: GetProduct :one
SELECT id, store_id, sku, name, description, price, stock_quantity, created_at, updated_at FROM products
WHERE id = $1 LIMIT 1
`

// Get a product by ID
func (q *Queries) GetProduct(ctx context.Context, id uuid.UUID) (Product, error) {
	row := q.db.QueryRowContext(ctx, getProduct, id)
	var i Product
	err := row.Scan(
		&i.ID,
		&i.StoreID,
		&i.Sku,
		&i.Name,
		&i.Description,
		&i.Price,
		&i.StockQuantity,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listProducts = `-- name: ListProducts :many
SELECT id, store_id, sku, name, description, price, stock_quantity, created_at, updated_at FROM products
ORDER BY name
`

// List all products
func (q *Queries) ListProducts(ctx context.Context) ([]Product, error) {
	rows, err := q.db.QueryContext(ctx, listProducts)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Product{}
	for rows.Next() {
		var i Product
		if err := rows.Scan(
			&i.ID,
			&i.StoreID,
			&i.Sku,
			&i.Name,
			&i.Description,
			&i.Price,
			&i.StockQuantity,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateProduct = `-- name: UpdateProduct :exec
UPDATE products
SET
  store_id = $2,
  sku = $3,
  name = $4,
  description = $5,
  price = $6,
  stock_quantity = $7,
  created_at = $8,
  updated_at = $9
WHERE id = $1
`

type UpdateProductParams struct {
	ID            uuid.UUID `json:"id"`
	StoreID       uuid.UUID `json:"storeId"`
	Sku           string    `json:"sku"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	Price         string    `json:"price"`
	StockQuantity int32     `json:"stockQuantity"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

// Update a product by ID
func (q *Queries) UpdateProduct(ctx context.Context, arg UpdateProductParams) error {
	_, err := q.db.ExecContext(ctx, updateProduct,
		arg.ID,
		arg.StoreID,
		arg.Sku,
		arg.Name,
		arg.Description,
		arg.Price,
		arg.StockQuantity,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	return err
}
