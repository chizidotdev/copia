package dto

import (
	"time"

	"github.com/google/uuid"
)

type DeleteOrderParams struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
}

type Order struct {
	ID                    uuid.UUID   `json:"id"`
	UserID                uuid.UUID   `json:"user_id" binding:"required"`
	CustomerID            uuid.UUID   `json:"customer_id"`
	Status                string      `json:"status" binding:"required"`
	EstimatedDeliveryDate time.Time   `json:"estimated_delivery_date" binding:"required"`
	OrderDate             time.Time   `json:"order_date" binding:"required"`
	TotalAmount           float32     `json:"total_amount" binding:"required"`
	OrderItems            []OrderItem `json:"order_items"`
}

type OrderItem struct {
	ID        uuid.UUID `json:"id"`
	OrderID   uuid.UUID `json:"order_id" binding:"required"`
	ProductID uuid.UUID `json:"product_id" binding:"required"`
	Quantity  int64     `json:"quantity" binding:"required"`
	UnitPrice float32   `json:"unit_price" binding:"required"`
	SubTotal  float32   `json:"sub_total" binding:"required"`
}

type UpdateOrderItemsRequest struct {
	OrderID    uuid.UUID `json:"order_id" binding:"required"`
	OrderItems []OrderItem
}
