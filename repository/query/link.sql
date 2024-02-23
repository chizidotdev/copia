-- Get a link by ID
-- name: GetLink :one
SELECT * FROM links
WHERE id = $1 LIMIT 1;

-- List all links
-- name: ListLinks :many
SELECT * FROM links
ORDER BY user_id;

-- Create a new link
-- name: CreateLink :one
INSERT INTO links (
  user_id, unique_link, link_type
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- Update a link by ID
-- name: UpdateLink :exec
UPDATE links
SET
  user_id = $2,
  unique_link = $3,
  link_type = $4
WHERE id = $1;

-- Delete a link by ID
-- name: DeleteLink :exec
DELETE FROM links
WHERE id = $1;

