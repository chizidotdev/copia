package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Queries interface {
	ListOrders(ctx context.Context, userEmail string) ([]Order, error)
	CreateOrder(ctx context.Context, arg CreateOrderParams) (Order, error)
	DeleteOrder(ctx context.Context, arg DeleteOrderParams) error
	GetOrder(ctx context.Context, id uuid.UUID) (Order, error)
	UpdateOrder(ctx context.Context, arg UpdateOrderParams) (Order, error)

	ListOrderItems(ctx context.Context, orderID uuid.UUID) ([]OrderItem, error)
	CreateOrderItem(ctx context.Context, arg CreateOrderItemParams) (OrderItem, error)
	DeleteOrderItem(ctx context.Context, id uuid.UUID) error
	GetOrderItem(ctx context.Context, id uuid.UUID) (OrderItem, error)
	UpdateOrderItem(ctx context.Context, arg UpdateOrderItemParams) (OrderItem, error)

	// GetReports(ctx context.Context, userEmail string) (datastruct.ReportRow, error)
}

var _ Queries = (*Store)(nil)

type Store struct {
	DB *gorm.DB
}

func NewStore(db *gorm.DB) *Store {
	err := db.AutoMigrate(
		&UserProfile{},
		&Customer{},
		&OrderItem{},
		&Order{},
		&Products{},
	)
	if err != nil {
		log.Panic("Cannot migrate db:", err)
	}

	return &Store{
		DB: db,
	}
}

func (s *Store) WithTx(tx *gorm.DB) *Store {
	return &Store{
		DB: tx,
	}
}

// ExecTx executes a function within a transaction.
func (s *Store) ExecTx(ctx context.Context, fn func(*Store) error) error {
	tx := s.DB.Begin()

	qtx := s.WithTx(tx)

	err := fn(qtx)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx error: %v, rb error: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit().Error
}
