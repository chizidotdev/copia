// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: store.sql

package repository

import (
	"context"

	"github.com/google/uuid"
)

const createStore = `-- name: CreateStore :one
INSERT INTO stores (
  user_id, name, description
) VALUES (
  $1, $2, $3
)
RETURNING id, user_id, name, description, image, created_at, updated_at
`

type CreateStoreParams struct {
	UserID      uuid.UUID `json:"userId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

// Create a new store
func (q *Queries) CreateStore(ctx context.Context, arg CreateStoreParams) (Store, error) {
	row := q.db.QueryRowContext(ctx, createStore, arg.UserID, arg.Name, arg.Description)
	var i Store
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Name,
		&i.Description,
		&i.Image,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteStore = `-- name: DeleteStore :exec
DELETE FROM stores
WHERE id = $1
`

// Delete a store by ID
func (q *Queries) DeleteStore(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteStore, id)
	return err
}

const getStore = `-- name: GetStore :one
SELECT id, user_id, name, description, image, created_at, updated_at FROM stores
WHERE id = $1 LIMIT 1
`

// Get a store by ID
func (q *Queries) GetStore(ctx context.Context, id uuid.UUID) (Store, error) {
	row := q.db.QueryRowContext(ctx, getStore, id)
	var i Store
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Name,
		&i.Description,
		&i.Image,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getStoreByUserId = `-- name: GetStoreByUserId :one
SELECT id, user_id, name, description, image, created_at, updated_at FROM stores
WHERE user_id = $1 LIMIT 1
`

// Get a store by user_id
func (q *Queries) GetStoreByUserId(ctx context.Context, userID uuid.UUID) (Store, error) {
	row := q.db.QueryRowContext(ctx, getStoreByUserId, userID)
	var i Store
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Name,
		&i.Description,
		&i.Image,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listStores = `-- name: ListStores :many
SELECT id, user_id, name, description, image, created_at, updated_at FROM stores
WHERE user_id = $1
ORDER BY name
`

// List all stores
func (q *Queries) ListStores(ctx context.Context, userID uuid.UUID) ([]Store, error) {
	rows, err := q.db.QueryContext(ctx, listStores, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Store{}
	for rows.Next() {
		var i Store
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Name,
			&i.Description,
			&i.Image,
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

const searchStores = `-- name: SearchStores :many
SELECT id, name FROM stores
WHERE name ILIKE '%' || $1::text || '%'
ORDER BY name
LIMIT 10
`

type SearchStoresRow struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// Search a store by name
func (q *Queries) SearchStores(ctx context.Context, query string) ([]SearchStoresRow, error) {
	rows, err := q.db.QueryContext(ctx, searchStores, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []SearchStoresRow{}
	for rows.Next() {
		var i SearchStoresRow
		if err := rows.Scan(&i.ID, &i.Name); err != nil {
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

const updateStore = `-- name: UpdateStore :many
UPDATE stores
SET
  name = $2,
  description = $3,
  updated_at = NOW()
WHERE id = $1
RETURNING id, user_id, name, description, image, created_at, updated_at
`

type UpdateStoreParams struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

// Update a store by ID
func (q *Queries) UpdateStore(ctx context.Context, arg UpdateStoreParams) ([]Store, error) {
	rows, err := q.db.QueryContext(ctx, updateStore, arg.ID, arg.Name, arg.Description)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Store{}
	for rows.Next() {
		var i Store
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Name,
			&i.Description,
			&i.Image,
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
