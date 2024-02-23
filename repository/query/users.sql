-- Get a user by ID
-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- List all users
-- name: ListUsers :many
SELECT * FROM users
ORDER BY first_name, last_name;

-- Create a new user
-- name: CreateUser :one
INSERT INTO users (
  email, first_name, last_name, image, password, google_id, role
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- Update a user by ID
-- name: UpdateUser :exec
UPDATE users
SET
  email = $2,
  first_name = $3,
  last_name = $4,
  image = $5,
  password = $6,
  google_id = $7,
  role = $8
WHERE id = $1;

-- Delete a user by ID
-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

