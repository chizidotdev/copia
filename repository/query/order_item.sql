-- Get an order item by ID
-- name: GetOrderItem :one
SELECT * FROM order_items
WHERE id = $1 LIMIT 1;

-- List all order items
-- name: ListOrderItems :many
SELECT * FROM order_items
ORDER BY order_id, product_id;

-- Create a new order item
-- name: CreateOrderItem :one
INSERT INTO order_items (
  order_id, product_id, quantity, unit_price, subtotal
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

-- Update an order item by ID
-- name: UpdateOrderItem :exec
UPDATE order_items
SET
  order_id = $2,
  product_id = $3,
  quantity = $4,
  unit_price = $5,
  subtotal = $6
WHERE id = $1;

-- Delete an order item by ID
-- name: DeleteOrderItem :exec
DELETE FROM order_items
WHERE id = $1;

