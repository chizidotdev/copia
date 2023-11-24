package usecases

import (
	"context"
	"fmt"

	"github.com/chizidotdev/copia/internal/app/core"
	"github.com/chizidotdev/copia/pkg/errors"
	"github.com/google/uuid"
)

type OrderRepository interface {
	ListOrders(ctx context.Context, userID uuid.UUID) ([]core.Order, error)
	GetOrder(ctx context.Context, id uuid.UUID) (core.Order, error)
	CreateOrder(ctx context.Context, arg core.Order) (core.Order, error)
	UpdateOrder(ctx context.Context, arg core.Order) (core.Order, error)
	UpdateOrderStatus(ctx context.Context, arg core.UpdateOrderStatusRequest) (core.Order, error)
	DeleteOrder(ctx context.Context, arg core.DeleteOrderRequest) error

	// OrderItems
	UpdateOrderItems(ctx context.Context, arg core.UpdateOrderItemsRequest) error
	DeleteOrderItems(ctx context.Context, arg core.DeleteOrderItemsRequest) error
}

type OrderService struct {
	Store OrderRepository
}

func NewOrderService(orderRepo OrderRepository) *OrderService {
	return &OrderService{
		Store: orderRepo,
	}
}

const (
	StatusPending   = "pending"
	StatusPaid      = "paid"
	StatusShipped   = "shipped"
	StatusDelivered = "delivered"
)

func validateOrderStatus(status string) error {
	switch status {
	case StatusPending, StatusPaid, StatusShipped, StatusDelivered:
		return nil
	default:
		msg := fmt.Sprintf(
			"Invalid status. Status must be one of %s, %s, %s, %s",
			StatusPending, StatusPaid, StatusShipped, StatusDelivered)
		return errors.Errorf(errors.ErrorBadRequest, msg)
	}
}

func (o *OrderService) CreateOrder(ctx context.Context, req core.OrderRequest) (core.Order, error) {
	if err := validateOrderStatus(req.Status); err != nil {
		return core.Order{}, err
	}

	totalAmount := float32(0)
	for _, orderItem := range req.OrderItems {
		totalAmount += orderItem.SubTotal
	}
	order, err := o.Store.CreateOrder(ctx, core.Order{
		UserID:                req.UserID,
		CustomerID:            req.CustomerID,
		Status:                req.Status,
		EstimatedDeliveryDate: req.EstimatedDeliveryDate,
		OrderDate:             req.OrderDate,
		TotalAmount:           totalAmount,
		OrderItems:            req.OrderItems,
	})
	if err != nil {
		return core.Order{}, errors.Errorf(errors.ErrorBadRequest, "Failed to create order")
	}

	return order, nil
}

func (o *OrderService) UpdateOrder(ctx context.Context, req core.UpdateOrderRequest) (core.Order, error) {
	order, err := o.Store.GetOrder(ctx, req.ID)
	if err != nil {
		return core.Order{}, errors.Errorf(errors.ErrorNotFound, "Order not found")
	}

	if err := validateOrderStatus(req.Status); err != nil {
		return core.Order{}, err
	}

	order.Status = req.Status
	order.EstimatedDeliveryDate = req.EstimatedDeliveryDate
	order.OrderDate = req.OrderDate
	order.OrderItems = req.OrderItems

	order, err = o.Store.UpdateOrder(ctx, order)
	if err != nil {
		return core.Order{}, errors.Errorf(errors.ErrorBadRequest, "Failed to update order. "+err.Error())
	}

	return order, nil
}

func (o *OrderService) UpdateOrderStatus(ctx context.Context, req core.UpdateOrderStatusRequest) (core.Order, error) {
	if err := validateOrderStatus(req.Status); err != nil {
		return core.Order{}, err
	}

	order, err := o.Store.UpdateOrderStatus(ctx, req)
	if err != nil {
		return core.Order{}, errors.Errorf(errors.ErrorBadRequest, "Failed to update order status. "+err.Error())
	}

	return order, nil
}

func (o *OrderService) ListOrders(ctx context.Context, userID uuid.UUID) ([]core.Order, error) {
	orders, err := o.Store.ListOrders(ctx, userID)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (o *OrderService) GetOrderByID(ctx context.Context, orderID uuid.UUID) (core.Order, error) {
	order, err := o.Store.GetOrder(ctx, orderID)
	if err != nil {
		return core.Order{}, errors.Errorf(errors.ErrorNotFound, "Order not found")
	}

	return order, nil
}

func (o *OrderService) DeleteOrder(ctx context.Context, req core.DeleteOrderRequest) error {
	err := o.Store.DeleteOrder(ctx, req)
	if err != nil {
		return errors.Errorf(errors.ErrorNotFound, "Order not found")
	}

	return nil
}

func (o *OrderService) UpdateOrderItems(ctx context.Context, req core.UpdateOrderItemsRequest) (core.Order, error) {
	order, err := o.Store.GetOrder(ctx, req.OrderID)
	if err != nil {
		return core.Order{}, errors.Errorf(errors.ErrorNotFound, "Order not found")
	}

	err = o.Store.UpdateOrderItems(ctx, req)
	if err != nil {
		return core.Order{}, errors.Errorf(errors.ErrorInternal, "Failed to update order items")
	}

	order.TotalAmount = 0
	for _, orderItem := range req.OrderItems {
		if orderItem.OrderID != req.OrderID {
			return core.Order{}, errors.Errorf(errors.ErrorBadRequest, "Order items must belong to the same order")
		}
		order.TotalAmount += orderItem.SubTotal
	}

	order, err = o.Store.UpdateOrder(ctx, order)
	if err != nil {
		return core.Order{}, errors.Errorf(errors.ErrorBadRequest, "Failed to update order")
	}

	return order, nil
}

func (o *OrderService) DeleteOrderItems(ctx context.Context, req core.DeleteOrderItemsRequest) (core.Order, error) {
	order, err := o.Store.GetOrder(ctx, req.OrderID)
	if err != nil {
		return core.Order{}, errors.Errorf(errors.ErrorNotFound, "Order not found")
	}

	err = o.Store.DeleteOrderItems(ctx, req)
	if err != nil {
		return core.Order{}, errors.Errorf(errors.ErrorInternal, "Failed to delete order items")
	}

	order.TotalAmount = 0
	for _, orderItem := range req.OrderItems {
		if orderItem.OrderID != req.OrderID {
			return core.Order{}, errors.Errorf(errors.ErrorBadRequest, "Order items must belong to the same order")
		}

		order.TotalAmount -= orderItem.SubTotal
	}

	order, err = o.Store.UpdateOrder(ctx, order)
	if err != nil {
		return core.Order{}, errors.Errorf(errors.ErrorInternal, "Failed to update order")
	}

	return order, nil
}
