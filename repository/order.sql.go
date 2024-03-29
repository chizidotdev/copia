// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: order.sql

package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createOrder = `-- name: CreateOrder :one
INSERT INTO orders (
  user_id, order_date, total_amount, payment_status, shipping_address, created_at, updated_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING id, user_id, order_date, total_amount, payment_status, shipping_address, created_at, updated_at
`

type CreateOrderParams struct {
	UserID          uuid.UUID     `json:"userId"`
	OrderDate       time.Time     `json:"orderDate"`
	TotalAmount     float64       `json:"totalAmount"`
	PaymentStatus   PaymentStatus `json:"paymentStatus"`
	ShippingAddress string        `json:"shippingAddress"`
	CreatedAt       time.Time     `json:"createdAt"`
	UpdatedAt       time.Time     `json:"updatedAt"`
}

// Create a new order
func (q *Queries) CreateOrder(ctx context.Context, arg CreateOrderParams) (Order, error) {
	row := q.db.QueryRowContext(ctx, createOrder,
		arg.UserID,
		arg.OrderDate,
		arg.TotalAmount,
		arg.PaymentStatus,
		arg.ShippingAddress,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i Order
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.OrderDate,
		&i.TotalAmount,
		&i.PaymentStatus,
		&i.ShippingAddress,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteOrder = `-- name: DeleteOrder :one
DELETE FROM orders
WHERE id = $1 AND user_id = $2
RETURNING id, user_id, order_date, total_amount, payment_status, shipping_address, created_at, updated_at
`

type DeleteOrderParams struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"userId"`
}

// Delete an order by ID
func (q *Queries) DeleteOrder(ctx context.Context, arg DeleteOrderParams) (Order, error) {
	row := q.db.QueryRowContext(ctx, deleteOrder, arg.ID, arg.UserID)
	var i Order
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.OrderDate,
		&i.TotalAmount,
		&i.PaymentStatus,
		&i.ShippingAddress,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getOrder = `-- name: GetOrder :one
SELECT id, user_id, order_date, total_amount, payment_status, shipping_address, created_at, updated_at FROM orders
WHERE id = $1 LIMIT 1
`

// Get an order by ID
func (q *Queries) GetOrder(ctx context.Context, id uuid.UUID) (Order, error) {
	row := q.db.QueryRowContext(ctx, getOrder, id)
	var i Order
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.OrderDate,
		&i.TotalAmount,
		&i.PaymentStatus,
		&i.ShippingAddress,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listUserOrders = `-- name: ListUserOrders :many
SELECT id, user_id, order_date, total_amount, payment_status, shipping_address, created_at, updated_at FROM orders
WHERE user_id = $1
ORDER BY order_date DESC
`

// List all orders for a user
func (q *Queries) ListUserOrders(ctx context.Context, userID uuid.UUID) ([]Order, error) {
	rows, err := q.db.QueryContext(ctx, listUserOrders, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Order{}
	for rows.Next() {
		var i Order
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.OrderDate,
			&i.TotalAmount,
			&i.PaymentStatus,
			&i.ShippingAddress,
			&i.CreatedAt,
			&i.UpdatedAt,
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

const updateOrder = `-- name: UpdateOrder :one
UPDATE orders
SET
  order_date = $3,
  total_amount = $4,
  payment_status = $5,
  shipping_address = $6,
  created_at = $7,
  updated_at = $8
WHERE id = $1 AND user_id = $2
RETURNING id, user_id, order_date, total_amount, payment_status, shipping_address, created_at, updated_at
`

type UpdateOrderParams struct {
	ID              uuid.UUID     `json:"id"`
	UserID          uuid.UUID     `json:"userId"`
	OrderDate       time.Time     `json:"orderDate"`
	TotalAmount     float64       `json:"totalAmount"`
	PaymentStatus   PaymentStatus `json:"paymentStatus"`
	ShippingAddress string        `json:"shippingAddress"`
	CreatedAt       time.Time     `json:"createdAt"`
	UpdatedAt       time.Time     `json:"updatedAt"`
}

// Update an order by ID
func (q *Queries) UpdateOrder(ctx context.Context, arg UpdateOrderParams) (Order, error) {
	row := q.db.QueryRowContext(ctx, updateOrder,
		arg.ID,
		arg.UserID,
		arg.OrderDate,
		arg.TotalAmount,
		arg.PaymentStatus,
		arg.ShippingAddress,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i Order
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.OrderDate,
		&i.TotalAmount,
		&i.PaymentStatus,
		&i.ShippingAddress,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
