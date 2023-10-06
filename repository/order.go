package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
)

func (s *Repository) ListOrders(_ context.Context, userEmail string) ([]Order, error) {
	var orders []Order
	result := s.DB.Preload("OrderItems").Find(&orders, "user_email = ?", userEmail)
	return orders, result.Error
}

func (s *Repository) GetOrder(_ context.Context, id uuid.UUID) (Order, error) {
	var order Order
	result := s.DB.Preload("OrderItems").First(&order, "id = ?", id)
	return order, result.Error
}

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

func (s *Repository) CreateOrder(_ context.Context, arg CreateOrderParams) (Order, error) {
	order := Order{
		Status:                arg.Status,
		ShippingDetails:       arg.ShippingDetails,
		EstimatedDeliveryDate: arg.EstimatedDeliveryDate,
		OrderDate:             arg.OrderDate,
		TotalAmount:           arg.TotalAmount,
		PaymentStatus:         arg.PaymentStatus,
		PaymentMethod:         arg.PaymentMethod,
		BillingAddress:        arg.BillingAddress,
		ShippingAddress:       arg.ShippingAddress,
		Notes:                 arg.Notes,
		CustomerID:            arg.CustomerID,
		UserEmail:             arg.UserEmail,
	}
	result := s.DB.Create(&order)
	return order, result.Error
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

func (s *Repository) UpdateOrder(_ context.Context, arg UpdateOrderParams) (Order, error) {
	var order Order
	if err := s.DB.First(&order, "id = ? AND user_email = ?", arg.ID, arg.UserEmail).Error; err != nil {
		return order, err
	}

	order.UserEmail = arg.UserEmail
	order.CustomerID = arg.CustomerID
	order.Status = arg.Status
	order.ShippingDetails = arg.ShippingDetails
	order.EstimatedDeliveryDate = arg.EstimatedDeliveryDate
	order.OrderDate = arg.OrderDate
	order.TotalAmount = arg.TotalAmount
	order.PaymentStatus = arg.PaymentStatus
	order.PaymentMethod = arg.PaymentMethod
	order.BillingAddress = arg.BillingAddress
	order.ShippingAddress = arg.ShippingAddress
	order.Notes = arg.Notes

	err := s.DB.Save(&order).Error
	return order, err
}

type DeleteOrderParams struct {
	ID        uuid.UUID `json:"id"`
	UserEmail string    `json:"user_email"`
}

func (s *Repository) DeleteOrder(_ context.Context, arg DeleteOrderParams) error {
	result := s.DB.Delete(&Order{}, "id = ? AND user_email = ?", arg.ID, arg.UserEmail)
	return result.Error
}
