-- Get a store by ID
-- name: GetStore :one
SELECT * FROM stores
WHERE id = $1 LIMIT 1;

-- List all stores
-- name: ListStores :many
SELECT * FROM stores
ORDER BY name;

-- Create a new store
-- name: CreateStore :one
INSERT INTO stores (
  user_id, name, description, created_at, updated_at
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

-- Update a store by ID
-- name: UpdateStore :exec
UPDATE stores
SET
  user_id = $2,
  name = $3,
  description = $4,
  created_at = $5,
  updated_at = $6
WHERE id = $1;

-- Delete a store by ID
-- name: DeleteStore :exec
DELETE FROM stores
WHERE id = $1;

