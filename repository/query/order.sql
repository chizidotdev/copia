-- Get an order by ID
-- name: GetOrder :one
SELECT * FROM orders
WHERE id = $1 LIMIT 1;

-- List all orders for a user
-- name: ListUserOrders :many
SELECT * FROM orders
WHERE user_id = $1
ORDER BY order_date;

-- Create a new order
-- name: CreateOrder :one
INSERT INTO orders (
  user_id, order_date, total_amount, payment_status, shipping_address, created_at, updated_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- Update an order by ID
-- name: UpdateOrder :one
UPDATE orders
SET
  order_date = $3,
  total_amount = $4,
  payment_status = $5,
  shipping_address = $6,
  created_at = $7,
  updated_at = $8
WHERE id = $1 AND user_id = $2
RETURNING *;

-- Delete an order by ID
-- name: DeleteOrder :one
DELETE FROM orders
WHERE id = $1 AND user_id = $2
RETURNING *;

