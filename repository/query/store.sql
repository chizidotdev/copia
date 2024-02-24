-- Get a store by ID
-- name: GetStore :one
SELECT * FROM stores
WHERE id = $1 LIMIT 1;

-- List all stores
-- name: ListStores :many
SELECT * FROM stores
WHERE user_id = $1
ORDER BY name;

-- Create a new store
-- name: CreateStore :one
INSERT INTO stores (
  user_id, name, description
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- Update a store by ID
-- name: UpdateStore :many
UPDATE stores
SET
  name = $2,
  description = $3,
  updated_at = NOW()
WHERE id = $1
RETURNING *;

-- Delete a store by ID
-- name: DeleteStore :exec
DELETE FROM stores
WHERE id = $1;

