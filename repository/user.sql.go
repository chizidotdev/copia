// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: user.sql

package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (
  email, first_name, last_name, image, google_id, role
) VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING id, email, first_name, last_name, image, google_id, role, created_at, updated_at
`

type CreateUserParams struct {
	Email     string         `json:"email"`
	FirstName string         `json:"firstName"`
	LastName  string         `json:"lastName"`
	Image     string         `json:"image"`
	GoogleID  sql.NullString `json:"googleId"`
	Role      UserRole       `json:"role"`
}

// Create a new user
func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.Email,
		arg.FirstName,
		arg.LastName,
		arg.Image,
		arg.GoogleID,
		arg.Role,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.FirstName,
		&i.LastName,
		&i.Image,
		&i.GoogleID,
		&i.Role,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1
`

// Delete a user by ID
func (q *Queries) DeleteUser(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteUser, id)
	return err
}

const getUser = `-- name: GetUser :one
SELECT id, email, first_name, last_name, image, google_id, role, created_at, updated_at FROM users
WHERE id = $1 LIMIT 1
`

// Get a user by ID
func (q *Queries) GetUser(ctx context.Context, id uuid.UUID) (User, error) {
	row := q.db.QueryRowContext(ctx, getUser, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.FirstName,
		&i.LastName,
		&i.Image,
		&i.GoogleID,
		&i.Role,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, email, first_name, last_name, image, google_id, role, created_at, updated_at FROM users
WHERE email = $1 LIMIT 1
`

// Get a user by email
func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.FirstName,
		&i.LastName,
		&i.Image,
		&i.GoogleID,
		&i.Role,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listUsers = `-- name: ListUsers :many
SELECT id, email, first_name, last_name, image, google_id, role, created_at, updated_at FROM users
ORDER BY first_name, last_name
`

// List all users
func (q *Queries) ListUsers(ctx context.Context) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, listUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []User{}
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Email,
			&i.FirstName,
			&i.LastName,
			&i.Image,
			&i.GoogleID,
			&i.Role,
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

const updateUser = `-- name: UpdateUser :exec
UPDATE users
SET
  email = $2,
  first_name = $3,
  last_name = $4,
  image = $5,
  google_id = $6,
  role = $7
WHERE id = $1
`

type UpdateUserParams struct {
	ID        uuid.UUID      `json:"id"`
	Email     string         `json:"email"`
	FirstName string         `json:"firstName"`
	LastName  string         `json:"lastName"`
	Image     string         `json:"image"`
	GoogleID  sql.NullString `json:"googleId"`
	Role      UserRole       `json:"role"`
}

// Update a user by ID
func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) error {
	_, err := q.db.ExecContext(ctx, updateUser,
		arg.ID,
		arg.Email,
		arg.FirstName,
		arg.LastName,
		arg.Image,
		arg.GoogleID,
		arg.Role,
	)
	return err
}

const upsertUser = `-- name: UpsertUser :one
INSERT INTO users (
  email, first_name, last_name, image, google_id, role
) VALUES (
  $1, $2, $3, $4, $5, $6
)
ON CONFLICT (email) DO UPDATE
SET
  email = $1,
  first_name = $2,
  last_name = $3,
  image = $4,
  google_id = $5,
  role = $6
RETURNING id, email, first_name, last_name, image, google_id, role, created_at, updated_at
`

type UpsertUserParams struct {
	Email     string         `json:"email"`
	FirstName string         `json:"firstName"`
	LastName  string         `json:"lastName"`
	Image     string         `json:"image"`
	GoogleID  sql.NullString `json:"googleId"`
	Role      UserRole       `json:"role"`
}

// Upsert a user by email
func (q *Queries) UpsertUser(ctx context.Context, arg UpsertUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, upsertUser,
		arg.Email,
		arg.FirstName,
		arg.LastName,
		arg.Image,
		arg.GoogleID,
		arg.Role,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.FirstName,
		&i.LastName,
		&i.Image,
		&i.GoogleID,
		&i.Role,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
