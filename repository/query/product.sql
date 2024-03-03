-- Get a product by ID
-- name: GetProduct :one
SELECT * FROM products
WHERE id = $1 LIMIT 1;

-- List all products
-- name: ListProducts :many
SELECT * FROM products
ORDER BY title;

-- List all products by store ID
-- name: ListProductsByStore :many
SELECT * FROM products
WHERE store_id = $1;

-- Create a new product
-- name: CreateProduct :one
INSERT INTO products (
  store_id, title, description, price, out_of_stock
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

-- Update a product by ID
-- name: UpdateProduct :one
UPDATE products
SET
  title = $3,
  description = $4,
  price = $5,
  out_of_stock = $6,
  updated_at = NOW()
WHERE id = $1 AND store_id = $2
RETURNING *;

-- Delete a product by ID
-- name: DeleteProduct :exec
DELETE FROM products
WHERE id = $1 AND store_id = $2;

