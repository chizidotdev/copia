// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: cart.sql

package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const clearCartItems = `-- name: ClearCartItems :exec
DELETE FROM cart_items
WHERE user_id = $1
`

// Clear cart items
func (q *Queries) ClearCartItems(ctx context.Context, userID uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, clearCartItems, userID)
	return err
}

const deleteCartItem = `-- name: DeleteCartItem :one
DELETE FROM cart_items
WHERE id = $1 AND user_id = $2
RETURNING id, user_id, product_id, quantity, created_at, updated_at
`

type DeleteCartItemParams struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"userId"`
}

// Remove a cart item
func (q *Queries) DeleteCartItem(ctx context.Context, arg DeleteCartItemParams) (CartItem, error) {
	row := q.db.QueryRowContext(ctx, deleteCartItem, arg.ID, arg.UserID)
	var i CartItem
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.ProductID,
		&i.Quantity,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getCartItems = `-- name: GetCartItems :many
SELECT
  ci.id, ci.user_id, ci.product_id, ci.quantity, ci.created_at, ci.updated_at,
  p.store_id,
  p.title,
  p.description,
  p.price,
  p.out_of_stock 
FROM
  cart_items ci
JOIN
  products p ON ci.product_id = p.id
WHERE
  ci.user_id = $1
`

type GetCartItemsRow struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"userId"`
	ProductID   uuid.UUID `json:"productId"`
	Quantity    int32     `json:"quantity"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	StoreID     uuid.UUID `json:"storeId"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	OutOfStock  bool      `json:"outOfStock"`
}

// Get cart items by user id
func (q *Queries) GetCartItems(ctx context.Context, userID uuid.UUID) ([]GetCartItemsRow, error) {
	rows, err := q.db.QueryContext(ctx, getCartItems, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetCartItemsRow{}
	for rows.Next() {
		var i GetCartItemsRow
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.ProductID,
			&i.Quantity,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.StoreID,
			&i.Title,
			&i.Description,
			&i.Price,
			&i.OutOfStock,
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

const updateCartItemQuantity = `-- name: UpdateCartItemQuantity :one
UPDATE cart_items
SET quantity = $2
WHERE id = $1 AND user_id = $3
RETURNING id, user_id, product_id, quantity, created_at, updated_at
`

type UpdateCartItemQuantityParams struct {
	ID       uuid.UUID `json:"id"`
	Quantity int32     `json:"quantity"`
	UserID   uuid.UUID `json:"userId"`
}

// Update cart item quantity
func (q *Queries) UpdateCartItemQuantity(ctx context.Context, arg UpdateCartItemQuantityParams) (CartItem, error) {
	row := q.db.QueryRowContext(ctx, updateCartItemQuantity, arg.ID, arg.Quantity, arg.UserID)
	var i CartItem
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.ProductID,
		&i.Quantity,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const upsertCartItem = `-- name: UpsertCartItem :one
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
RETURNING id, user_id, product_id, quantity, created_at, updated_at
`

type UpsertCartItemParams struct {
	UserID    uuid.UUID `json:"userId"`
	ProductID uuid.UUID `json:"productId"`
	Quantity  int32     `json:"quantity"`
}

// Add or update a cart item
func (q *Queries) UpsertCartItem(ctx context.Context, arg UpsertCartItemParams) (CartItem, error) {
	row := q.db.QueryRowContext(ctx, upsertCartItem, arg.UserID, arg.ProductID, arg.Quantity)
	var i CartItem
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.ProductID,
		&i.Quantity,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
