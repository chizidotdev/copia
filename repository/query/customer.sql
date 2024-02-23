-- Get a customer by ID
-- name: GetCustomer :one
SELECT * FROM customers
WHERE id = $1 LIMIT 1;

-- List all customers
-- name: ListCustomers :many
SELECT * FROM customers
ORDER BY store_id, first_name, last_name;

-- Create a new customer
-- name: CreateCustomer :one
INSERT INTO customers (
  store_id, first_name, last_name, email, phone, address
) VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- Update a customer by ID
-- name: UpdateCustomer :exec
UPDATE customers
SET
  store_id = $2,
  first_name = $3,
  last_name = $4,
  email = $5,
  phone = $6,
  address = $7
WHERE id = $1;

-- Delete a customer by ID
-- name: DeleteCustomer :exec
DELETE FROM customers
WHERE id = $1;

