// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: order_item.sql

package repository

import (
	"context"

	"github.com/google/uuid"
)

const createOrderItem = `-- name: CreateOrderItem :one
INSERT INTO order_items (
  order_id, product_id, quantity, unit_price, subtotal
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING id, order_id, product_id, quantity, unit_price, subtotal
`

type CreateOrderItemParams struct {
	OrderID   uuid.UUID `json:"orderId"`
	ProductID uuid.UUID `json:"productId"`
	Quantity  int32     `json:"quantity"`
	UnitPrice float64   `json:"unitPrice"`
	Subtotal  float64   `json:"subtotal"`
}

// Create a new order item
func (q *Queries) CreateOrderItem(ctx context.Context, arg CreateOrderItemParams) (OrderItem, error) {
	row := q.db.QueryRowContext(ctx, createOrderItem,
		arg.OrderID,
		arg.ProductID,
		arg.Quantity,
		arg.UnitPrice,
		arg.Subtotal,
	)
	var i OrderItem
	err := row.Scan(
		&i.ID,
		&i.OrderID,
		&i.ProductID,
		&i.Quantity,
		&i.UnitPrice,
		&i.Subtotal,
	)
	return i, err
}

const deleteOrderItem = `-- name: DeleteOrderItem :exec
DELETE FROM order_items
WHERE id = $1
`

// Delete an order item by ID
func (q *Queries) DeleteOrderItem(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteOrderItem, id)
	return err
}

const getOrderItem = `-- name: GetOrderItem :one
SELECT id, order_id, product_id, quantity, unit_price, subtotal FROM order_items
WHERE id = $1 LIMIT 1
`

// Get an order item by ID
func (q *Queries) GetOrderItem(ctx context.Context, id uuid.UUID) (OrderItem, error) {
	row := q.db.QueryRowContext(ctx, getOrderItem, id)
	var i OrderItem
	err := row.Scan(
		&i.ID,
		&i.OrderID,
		&i.ProductID,
		&i.Quantity,
		&i.UnitPrice,
		&i.Subtotal,
	)
	return i, err
}

const listOrderItems = `-- name: ListOrderItems :many
SELECT id, order_id, product_id, quantity, unit_price, subtotal FROM order_items
ORDER BY order_id, product_id
`

// List all order items
func (q *Queries) ListOrderItems(ctx context.Context) ([]OrderItem, error) {
	rows, err := q.db.QueryContext(ctx, listOrderItems)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []OrderItem{}
	for rows.Next() {
		var i OrderItem
		if err := rows.Scan(
			&i.ID,
			&i.OrderID,
			&i.ProductID,
			&i.Quantity,
			&i.UnitPrice,
			&i.Subtotal,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateOrderItem = `-- name: UpdateOrderItem :exec
UPDATE order_items
SET
  order_id = $2,
  product_id = $3,
  quantity = $4,
  unit_price = $5,
  subtotal = $6
WHERE id = $1
`

type UpdateOrderItemParams struct {
	ID        uuid.UUID `json:"id"`
	OrderID   uuid.UUID `json:"orderId"`
	ProductID uuid.UUID `json:"productId"`
	Quantity  int32     `json:"quantity"`
	UnitPrice float64   `json:"unitPrice"`
	Subtotal  float64   `json:"subtotal"`
}

// Update an order item by ID
func (q *Queries) UpdateOrderItem(ctx context.Context, arg UpdateOrderItemParams) error {
	_, err := q.db.ExecContext(ctx, updateOrderItem,
		arg.ID,
		arg.OrderID,
		arg.ProductID,
		arg.Quantity,
		arg.UnitPrice,
		arg.Subtotal,
	)
	return err
}
