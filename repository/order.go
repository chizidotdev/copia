package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
)

func (r *Repository) ListOrders(_ context.Context, userID uuid.UUID) ([]Order, error) {
	var orders []Order
	result := r.DB.Preload("OrderItems").Find(&orders, "user_id = ?", userID)
	return orders, result.Error
}

func (r *Repository) GetOrder(_ context.Context, id uuid.UUID) (Order, error) {
	var order Order
	result := r.DB.Preload("OrderItems").First(&order, "id = ?", id)
	return order, result.Error
}

type CreateOrderParams struct {
	UserID                uuid.UUID `json:"user_id"`
	CustomerID            uuid.UUID `json:"customer_id"`
	Status                string    `json:"status"`
	EstimatedDeliveryDate time.Time `json:"estimated_delivery_date"`
	OrderDate             time.Time `json:"order_date"`
	TotalAmount           float32   `json:"total_amount"`
}

func (r *Repository) CreateOrder(_ context.Context, arg CreateOrderParams) (Order, error) {
	order := Order{
		Status:                arg.Status,
		EstimatedDeliveryDate: arg.EstimatedDeliveryDate,
		OrderDate:             arg.OrderDate,
		TotalAmount:           arg.TotalAmount,
		CustomerID:            arg.CustomerID,
		UserID:                arg.UserID,
	}
	result := r.DB.Create(&order)
	return order, result.Error
}

type UpdateOrderParams struct {
	ID                    uuid.UUID `json:"id"`
	UserID                uuid.UUID `json:"user_id"`
	CustomerID            uuid.UUID `json:"customer_id"`
	Status                string    `json:"status"`
	EstimatedDeliveryDate time.Time `json:"estimated_delivery_date"`
	OrderDate             time.Time `json:"order_date"`
	TotalAmount           float32   `json:"total_amount"`
}

func (r *Repository) UpdateOrder(_ context.Context, arg UpdateOrderParams) (Order, error) {
	var order Order
	if err := r.DB.First(&order, "id = ? AND user_id = ?", arg.ID, arg.UserID).Error; err != nil {
		return order, err
	}

	order.UserID = arg.UserID
	order.CustomerID = arg.CustomerID
	order.Status = arg.Status
	order.EstimatedDeliveryDate = arg.EstimatedDeliveryDate
	order.OrderDate = arg.OrderDate
	order.TotalAmount = arg.TotalAmount

	err := r.DB.Save(&order).Error
	return order, err
}

type DeleteOrderParams struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
}

func (r *Repository) DeleteOrder(_ context.Context, arg DeleteOrderParams) error {
	result := r.DB.Delete(&Order{}, "id = ? AND user_id = ?", arg.ID, arg.UserID)
	return result.Error
}
