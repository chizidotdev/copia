package repository

import (
	"context"

	"github.com/chizidotdev/copia/internal/datastruct"
	"github.com/google/uuid"
)

func (s *Store) ListOrders(_ context.Context, userEmail string) ([]Order, error) {
	var orders []Order
	result := s.DB.Find(&orders, "user_email = ?", userEmail)
	return orders, result.Error
}

func (s *Store) GetOrder(_ context.Context, id uuid.UUID) (Order, error) {
	var order Order
	result := s.DB.First(&order, "id = ?", id)
	return order, result.Error
}

func (s *Store) CreateOrder(_ context.Context, arg datastruct.CreateOrderParams) (Order, error) {
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

func (s *Store) UpdateOrder(_ context.Context, arg datastruct.UpdateOrderParams) (Order, error) {
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

func (s *Store) DeleteOrder(_ context.Context, arg datastruct.DeleteOrderParams) error {
	result := s.DB.Delete(&Order{}, "id = ? AND user_email = ?", arg.ID, arg.UserEmail)
	return result.Error
}
