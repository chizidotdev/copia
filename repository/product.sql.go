// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: product.sql

package repository

import (
	"context"

	"github.com/google/uuid"
)

const createProduct = `-- name: CreateProduct :one
INSERT INTO products (
  store_id, title, description, price, out_of_stock
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING id, store_id, title, description, price, out_of_stock, created_at, updated_at, deleted_at
`

type CreateProductParams struct {
	StoreID     uuid.UUID `json:"storeId"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	OutOfStock  bool      `json:"outOfStock"`
}

// Create a new product
func (q *Queries) CreateProduct(ctx context.Context, arg CreateProductParams) (Product, error) {
	row := q.db.QueryRowContext(ctx, createProduct,
		arg.StoreID,
		arg.Title,
		arg.Description,
		arg.Price,
		arg.OutOfStock,
	)
	var i Product
	err := row.Scan(
		&i.ID,
		&i.StoreID,
		&i.Title,
		&i.Description,
		&i.Price,
		&i.OutOfStock,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const deleteProduct = `-- name: DeleteProduct :exec
UPDATE products
SET deleted_at = NOW()
WHERE id = $1 AND store_id = $2
`

type DeleteProductParams struct {
	ID      uuid.UUID `json:"id"`
	StoreID uuid.UUID `json:"storeId"`
}

// Delete a product by ID
func (q *Queries) DeleteProduct(ctx context.Context, arg DeleteProductParams) error {
	_, err := q.db.ExecContext(ctx, deleteProduct, arg.ID, arg.StoreID)
	return err
}

const getProduct = `-- name: GetProduct :one
SELECT id, store_id, title, description, price, out_of_stock, created_at, updated_at, deleted_at FROM products
WHERE id = $1 
AND deleted_at IS NULL 
LIMIT 1
`

// Get a product by ID
func (q *Queries) GetProduct(ctx context.Context, id uuid.UUID) (Product, error) {
	row := q.db.QueryRowContext(ctx, getProduct, id)
	var i Product
	err := row.Scan(
		&i.ID,
		&i.StoreID,
		&i.Title,
		&i.Description,
		&i.Price,
		&i.OutOfStock,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const listProductsByStore = `-- name: ListProductsByStore :many
SELECT id, store_id, title, description, price, out_of_stock, created_at, updated_at, deleted_at FROM products
WHERE store_id = $1
AND deleted_at IS NULL
ORDER BY created_at ASC
`

// List all products by store ID
func (q *Queries) ListProductsByStore(ctx context.Context, storeID uuid.UUID) ([]Product, error) {
	rows, err := q.db.QueryContext(ctx, listProductsByStore, storeID)
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
			&i.Title,
			&i.Description,
			&i.Price,
			&i.OutOfStock,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
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

const searchProducts = `-- name: SearchProducts :many
SELECT id, store_id, title, description, price, out_of_stock, created_at, updated_at, deleted_at FROM products
WHERE (title ILIKE '%' || $1::text || '%' OR description ILIKE '%' || $1::text || '%')
AND deleted_at IS NULL
ORDER BY title
LIMIT 10
`

// Search a stores product by title and description
func (q *Queries) SearchProducts(ctx context.Context, query string) ([]Product, error) {
	rows, err := q.db.QueryContext(ctx, searchProducts, query)
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
			&i.Title,
			&i.Description,
			&i.Price,
			&i.OutOfStock,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
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

const updateProduct = `-- name: UpdateProduct :one
UPDATE products
SET
  title = $3,
  description = $4,
  price = $5,
  out_of_stock = $6,
  updated_at = NOW()
WHERE id = $1 AND store_id = $2
RETURNING id, store_id, title, description, price, out_of_stock, created_at, updated_at, deleted_at
`

type UpdateProductParams struct {
	ID          uuid.UUID `json:"id"`
	StoreID     uuid.UUID `json:"storeId"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	OutOfStock  bool      `json:"outOfStock"`
}

// Update a product by ID
func (q *Queries) UpdateProduct(ctx context.Context, arg UpdateProductParams) (Product, error) {
	row := q.db.QueryRowContext(ctx, updateProduct,
		arg.ID,
		arg.StoreID,
		arg.Title,
		arg.Description,
		arg.Price,
		arg.OutOfStock,
	)
	var i Product
	err := row.Scan(
		&i.ID,
		&i.StoreID,
		&i.Title,
		&i.Description,
		&i.Price,
		&i.OutOfStock,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}
