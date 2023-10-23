package adapters

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// Base contains common columns for all tables.
type Base struct {
	ID        uuid.UUID      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (base *Base) BeforeCreate(db *gorm.DB) error {
	base.ID = uuid.New()

	return nil
}

type Repository[R any] struct {
	Repo *R
}

func NewRepository[R any](repo *R) *Repository[R] {
	return &Repository[R]{
		Repo: repo,
	}
}

func (r *Repository[R]) makeExecTx(db *gorm.DB) func(ctx context.Context, fn func(*R) error) error {
	return func(ctx context.Context, fn func(*R) error) error {
		tx := db.Begin()

		err := fn(r.Repo)
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				return fmt.Errorf("tx errors: %v, rb errors: %v", err, rbErr)
			}
			return err
		}
		return tx.Commit().Error
	}
}
