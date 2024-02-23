-- Get a commission by ID
-- name: GetCommission :one
SELECT * FROM commissions
WHERE id = $1 LIMIT 1;

-- List all commissions
-- name: ListCommissions :many
SELECT * FROM commissions
ORDER BY order_id, user_id;

-- Create a new commission
-- name: CreateCommission :one
INSERT INTO commissions (
  order_id, user_id, commission_amount, paid_status
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- Update a commission by ID
-- name: UpdateCommission :exec
UPDATE commissions
SET
  order_id = $2,
  user_id = $3,
  commission_amount = $4,
  paid_status = $5
WHERE id = $1;

-- Delete a commission by ID
-- name: DeleteCommission :exec
DELETE FROM commissions
WHERE id = $1;

