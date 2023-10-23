package core

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID                    uuid.UUID   `json:"id"`
	UserID                uuid.UUID   `json:"user_id"`
	CustomerID            uuid.UUID   `json:"customer_id"`
	Status                string      `json:"status"`
	EstimatedDeliveryDate time.Time   `json:"estimated_delivery_date"`
	OrderDate             time.Time   `json:"order_date"`
	TotalAmount           float32     `json:"total_amount"`
	OrderItems            []OrderItem `json:"order_item"`
}

type OrderRequest struct {
	UserID                uuid.UUID   `json:"user_id"`
	CustomerID            uuid.UUID   `json:"customer_id"`
	Status                string      `json:"status" binding:"required"`
	EstimatedDeliveryDate time.Time   `json:"estimated_delivery_date" binding:"required"`
	OrderDate             time.Time   `json:"order_date" binding:"required"`
	OrderItems            []OrderItem `json:"order_items"`
}

type DeleteOrderRequest struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
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
	OrderID    uuid.UUID   `json:"order_id" binding:"required"`
	OrderItems []OrderItem `json:"order_items"`
}
