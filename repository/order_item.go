package repository

import (
	"context"

	"github.com/google/uuid"
)

func (r *Repository) ListOrderItems(_ context.Context, orderID uuid.UUID) ([]OrderItem, error) {
	var orderItems []OrderItem
	result := r.DB.Find(&orderItems, "order_id = ?", orderID)
	return orderItems, result.Error
}

func (r *Repository) GetOrderItem(_ context.Context, id uuid.UUID) (OrderItem, error) {
	var orderItem OrderItem
	result := r.DB.First(&orderItem, "id = ?", id)
	return orderItem, result.Error
}

type CreateOrderItemParams struct {
	OrderID   uuid.UUID `json:"order_id"`
	ProductID uuid.UUID `json:"product_id"`
	Quantity  int64     `json:"quantity"`
	UnitPrice float32   `json:"unit_price"`
	SubTotal  float32   `json:"sub_total"`
}

func (r *Repository) CreateOrderItem(_ context.Context, arg CreateOrderItemParams) (OrderItem, error) {
	orderItem := OrderItem{
		OrderID:   arg.OrderID,
		ProductID: arg.ProductID,
		Quantity:  arg.Quantity,
		UnitPrice: arg.UnitPrice,
		SubTotal:  arg.SubTotal,
	}
	result := r.DB.Create(&orderItem)
	return orderItem, result.Error
}

type UpdateOrderItemParams struct {
	ID        uuid.UUID `json:"id"`
	OrderID   uuid.UUID `json:"order_id"`
	ProductID uuid.UUID `json:"product_id"`
	Quantity  int64     `json:"quantity"`
	UnitPrice float32   `json:"unit_price"`
	SubTotal  float32   `json:"sub_total"`
}

func (r *Repository) UpdateOrderItem(_ context.Context, arg UpdateOrderItemParams) (OrderItem, error) {
	var orderItem OrderItem
	if err := r.DB.First(&orderItem, "id = ?", arg.ID).Error; err != nil {
		return orderItem, err
	}

	orderItem.OrderID = arg.OrderID
	orderItem.ProductID = arg.ProductID
	orderItem.Quantity = arg.Quantity
	orderItem.UnitPrice = arg.UnitPrice
	orderItem.SubTotal = arg.SubTotal

	err := r.DB.Save(&orderItem).Error
	return orderItem, err
}

func (r *Repository) DeleteOrderItem(_ context.Context, id uuid.UUID) error {
	result := r.DB.Delete(&OrderItem{}, "id = ?", id)
	return result.Error
}
