-- Get a customer by ID
-- name: GetCustomer :one
SELECT c.*,
  u.first_name,
  u.last_name,
  u.email,
  u.image
FROM customers c
JOIN users u ON c.user_id = u.id
WHERE c.id = $1 AND c.store_id = $2
LIMIT 1;

-- List all customers for a store
-- name: ListCustomers :many
SELECT c.*,
  u.first_name,
  u.last_name,
  u.email,
  u.image
FROM customers c
JOIN users u ON c.user_id = u.id
WHERE c.store_id = $1
ORDER BY c.created_at;

-- Create a new customer
-- name: CreateCustomer :one
INSERT INTO customers (
  store_id, user_id
) VALUES (
  $1, $2
)
RETURNING *;

-- Delete a customer by ID
-- name: DeleteCustomer :exec
DELETE FROM customers
WHERE id = $1 AND store_id = $2;

