-- Get a product image by ID:
-- name: GetProductImage :one
SELECT * FROM product_images
WHERE id = $1 LIMIT 1;

-- List all product images for a product:
-- name: ListProductImages :many
SELECT * FROM product_images
WHERE product_id = $1;

-- Create a product image for a product:
-- name: CreateProductImage :one
INSERT INTO product_images (product_id, url)
VALUES ($1, $2)
RETURNING *;

-- Create product images for a product:
-- name: BulkCreateProductImages :many
INSERT INTO product_images (product_id, url)
SELECT 
  unnest(sqlc.arg(product_ids)::uuid[]),
  unnest(sqlc.arg(urls)::varchar[])
RETURNING *;

-- Update a product image by ID:
-- name: UpdateProductImageURL :exec
UPDATE product_images
SET
  url = $2
WHERE id = $1
RETURNING *;

-- Set a product image as primary and others as non-primary:
-- name: SetPrimaryImage :exec
UPDATE product_images
SET is_primary = CASE
  WHEN id = $1 THEN true
  ELSE false
END
WHERE product_id = $2;

-- Delete a product image by ID:
-- name: DeleteProductImage :exec
DELETE FROM product_images
WHERE id = $1
RETURNING *;
