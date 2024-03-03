package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type Repository struct {
	*Queries
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Queries: New(db),
		db:      db,
	}
}

func (r *Repository) ExecTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	queries := New(tx)
	err = fn(queries)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx error: %v, rb error: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

func ParseUUID(value string) (uuid.UUID, error) {
	return uuid.Parse(value)
}
