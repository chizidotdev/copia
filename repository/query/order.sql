-- Get an order by ID
-- name: GetOrder :one
SELECT * FROM orders
WHERE id = $1 LIMIT 1;

-- List all orders
-- name: ListOrders :many
SELECT * FROM orders
ORDER BY order_date;

-- Create a new order
-- name: CreateOrder :one
INSERT INTO orders (
  user_id, order_date, total_amount, status, payment_status, shipping_address, created_at, updated_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- Update an order by ID
-- name: UpdateOrder :exec
UPDATE orders
SET
  user_id = $2,
  order_date = $3,
  total_amount = $4,
  status = $5,
  payment_status = $6,
  shipping_address = $7,
  created_at = $8,
  updated_at = $9
WHERE id = $1;

-- Delete an order by ID
-- name: DeleteOrder :exec
DELETE FROM orders
WHERE id = $1;

