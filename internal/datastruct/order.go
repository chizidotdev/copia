package datastruct

import (
	"time"

	"github.com/google/uuid"
)

type CreateOrderParams struct {
	UserEmail             string    `json:"user_email"`
	CustomerID            uuid.UUID `json:"customer_id"`
	Status                string    `json:"status"`
	ShippingDetails       string    `json:"shipping_details"`
	EstimatedDeliveryDate time.Time `json:"estimated_delivery_date"`
	OrderDate             time.Time `json:"order_date"`
	TotalAmount           float32   `json:"total_amount"`
	PaymentStatus         string    `json:"payment_status"`
	PaymentMethod         string    `json:"payment_method"`
	BillingAddress        string    `json:"billing_address"`
	ShippingAddress       string    `json:"shipping_address"`
	Notes                 string    `json:"notes"`
}

type DeleteOrderParams struct {
	ID        uuid.UUID `json:"id"`
	UserEmail string    `json:"user_email"`
}

type UpdateOrderParams struct {
	ID                    uuid.UUID `json:"id"`
	UserEmail             string    `json:"user_email"`
	CustomerID            uuid.UUID `json:"customer_id"`
	Status                string    `json:"status"`
	ShippingDetails       string    `json:"shipping_details"`
	EstimatedDeliveryDate time.Time `json:"estimated_delivery_date"`
	OrderDate             time.Time `json:"order_date"`
	TotalAmount           float32   `json:"total_amount"`
	PaymentStatus         string    `json:"payment_status"`
	PaymentMethod         string    `json:"payment_method"`
	BillingAddress        string    `json:"billing_address"`
	ShippingAddress       string    `json:"shipping_address"`
	Notes                 string    `json:"notes"`
}
