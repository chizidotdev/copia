package datastruct

import "github.com/google/uuid"

type CreateOrderItemParams struct {
	OrderID   uuid.UUID `json:"order_id"`
	ProductID uuid.UUID `json:"product_id"`
	Quantity  int64     `json:"quantity"`
	UnitPrice float32   `json:"unit_price"`
	SubTotal  float32   `json:"sub_total"`
}

type UpdateOrderItemParams struct {
	ID        uuid.UUID `json:"id"`
	OrderID   uuid.UUID `json:"order_id"`
	ProductID uuid.UUID `json:"product_id"`
	Quantity  int64     `json:"quantity"`
	UnitPrice float32   `json:"unit_price"`
	SubTotal  float32   `json:"sub_total"`
}
