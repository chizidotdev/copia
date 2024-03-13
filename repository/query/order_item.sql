-- Get an order item by ID
-- name: GetOrderItem :one
SELECT * FROM order_items
WHERE id = $1 LIMIT 1;

-- List all order items for an order
-- name: ListOrderItems :many
SELECT oi.* ,
  p.title AS product_title
FROM order_items oi
JOIN products p ON oi.product_id = p.id
WHERE oi.order_id = $1
ORDER BY oi.order_id, oi.product_id;

-- List all order items for a store
-- name: ListStoreOrderItems :many
SELECT oi.*,
  p.title AS product_title, 
  o.order_date AS order_date,
  o.payment_status AS payment_status,
  o.shipping_address AS shipping_address
FROM order_items oi
JOIN products p ON oi.product_id = p.id
JOIN orders o ON oi.order_id = o.id
WHERE oi.store_id = $1
ORDER BY o.order_date DESC;

-- Create a new order item
-- name: CreateOrderItem :one
INSERT INTO order_items (
  order_id, product_id, store_id, status, quantity, unit_price, subtotal
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- Update an order item status
-- name: UpdateOrderItem :one
UPDATE order_items
SET status = $3
WHERE id = $1 AND store_id = $2
RETURNING *;

-- Delete an order item by ID
-- name: DeleteOrderItem :one
DELETE FROM order_items
WHERE id = $1
RETURNING *;
