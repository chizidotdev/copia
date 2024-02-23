-- Get a product by ID
-- name: GetProduct :one
SELECT * FROM products
WHERE id = $1 LIMIT 1;

-- List all products
-- name: ListProducts :many
SELECT * FROM products
ORDER BY name;

-- Create a new product
-- name: CreateProduct :one
INSERT INTO products (
  store_id, sku, name, description, price, stock_quantity, created_at, updated_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- Update a product by ID
-- name: UpdateProduct :exec
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
WHERE id = $1;

-- Delete a product by ID
-- name: DeleteProduct :exec
DELETE FROM products
WHERE id = $1;

