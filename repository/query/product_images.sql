-- Get a product image by ID:
-- name: GetProductImage :one
SELECT * FROM product_images
WHERE id = $1 LIMIT 1;

-- List all product images for a product:
-- name: ListProductImagesForProduct :many
SELECT * FROM product_images
WHERE product_id = $1;

-- Create a product image for a product:
-- name: CreateProductImage :one
INSERT INTO product_images (product_id, url)
VALUES ($1, $2)
RETURNING *;

-- Update a product image by ID:
-- name: UpdateProductImage :exec
UPDATE product_images
SET
  url = $2
WHERE id = $1
RETURNING *;

-- Delete a product image by ID:
-- name: DeleteProductImage :exec
DELETE FROM product_images
WHERE id = $1
RETURNING *;
