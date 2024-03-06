-- Add or update a cart item
-- name: UpsertCartItem :one
INSERT INTO cart_items (
    user_id,
    product_id,
    quantity
) VALUES (
    $1,
    $2,
    $3
) ON CONFLICT (user_id, product_id) DO UPDATE
SET quantity = cart_items.quantity + $3
RETURNING *;

-- Get cart items by user id
-- name: GetCartItems :many
SELECT
  ci.*,
  p.title,
  p.description,
  p.price,
  p.out_of_stock
FROM
  cart_items ci
JOIN
  products p ON ci.product_id = p.id
WHERE
  ci.user_id = $1;

-- Remove a cart item
-- name: DeleteCartItem :one
DELETE FROM cart_items
WHERE id = $1 AND user_id = $2
RETURNING *;

-- Update cart item quantity
-- name: UpdateCartItemQuantity :one
UPDATE cart_items
SET quantity = $2
WHERE id = $1 AND user_id = $3
RETURNING *;
