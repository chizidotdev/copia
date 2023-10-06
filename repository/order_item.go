package repository

import (
	"context"

	"github.com/google/uuid"
)

func (s *Repository) ListOrderItems(_ context.Context, orderID uuid.UUID) ([]OrderItem, error) {
	var orderItems []OrderItem
	result := s.DB.Find(&orderItems, "order_id = ?", orderID)
	return orderItems, result.Error
}

func (s *Repository) GetOrderItem(_ context.Context, id uuid.UUID) (OrderItem, error) {
	var orderItem OrderItem
	result := s.DB.First(&orderItem, "id = ?", id)
	return orderItem, result.Error
}

type CreateOrderItemParams struct {
	OrderID   uuid.UUID `json:"order_id"`
	ProductID uuid.UUID `json:"product_id"`
	Quantity  int64     `json:"quantity"`
	UnitPrice float32   `json:"unit_price"`
	SubTotal  float32   `json:"sub_total"`
}

func (s *Repository) CreateOrderItem(_ context.Context, arg CreateOrderItemParams) (OrderItem, error) {
	orderItem := OrderItem{
		OrderID:   arg.OrderID,
		ProductID: arg.ProductID,
		Quantity:  arg.Quantity,
		UnitPrice: arg.UnitPrice,
		SubTotal:  arg.SubTotal,
	}
	result := s.DB.Create(&orderItem)
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

func (s *Repository) UpdateOrderItem(_ context.Context, arg UpdateOrderItemParams) (OrderItem, error) {
	var orderItem OrderItem
	if err := s.DB.First(&orderItem, "id = ?", arg.ID).Error; err != nil {
		return orderItem, err
	}

	orderItem.OrderID = arg.OrderID
	orderItem.ProductID = arg.ProductID
	orderItem.Quantity = arg.Quantity
	orderItem.UnitPrice = arg.UnitPrice
	orderItem.SubTotal = arg.SubTotal

	err := s.DB.Save(&orderItem).Error
	return orderItem, err
}

func (s *Repository) DeleteOrderItem(_ context.Context, id uuid.UUID) error {
	result := s.DB.Delete(&OrderItem{}, "id = ?", id)
	return result.Error
}
