package usecases

import (
	"context"
	"github.com/chizidotdev/copia/internal/app/core"
	"github.com/chizidotdev/copia/pkg/errors"
	"github.com/google/uuid"
)

type OrderRepository interface {
	ListOrders(ctx context.Context, userID uuid.UUID) ([]core.Order, error)
	GetOrder(ctx context.Context, id uuid.UUID) (core.Order, error)
	CreateOrder(ctx context.Context, arg core.Order) (core.Order, error)
	DeleteOrder(ctx context.Context, arg core.DeleteOrderRequest) error
}

type OrderService struct {
	Store OrderRepository
}

func NewOrderService(orderRepo OrderRepository) *OrderService {
	return &OrderService{
		Store: orderRepo,
	}
}

func (o *OrderService) CreateOrder(ctx context.Context, req core.OrderRequest) (core.Order, error) {
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
		errResp := errors.ErrResponse{
			Code:      errors.ErrorBadRequest,
			MessageID: "",
			Message:   "Failed to create order",
			Reason:    err.Error(),
		}
		return core.Order{}, errors.Errorf(errResp)
	}

	return order, nil
}

func (o *OrderService) ListOrders(ctx context.Context, userID uuid.UUID) ([]core.Order, error) {
	orders, err := o.Store.ListOrders(ctx, userID)
	if err != nil {
		errResp := errors.ErrResponse{
			Code:      errors.ErrorBadRequest,
			MessageID: "",
			Message:   "Failed to get orders",
			Reason:    err.Error(),
		}
		return nil, errors.Errorf(errResp)
	}

	return orders, nil
}

func (o *OrderService) GetOrderByID(ctx context.Context, orderID uuid.UUID) (core.Order, error) {
	order, err := o.Store.GetOrder(ctx, orderID)
	if err != nil {
		errResp := errors.ErrResponse{
			Code:      errors.ErrorNotFound,
			MessageID: "",
			Message:   "Order not found",
			Reason:    err.Error(),
		}
		return core.Order{}, errors.Errorf(errResp)
	}

	return order, nil
}

func (o *OrderService) DeleteOrder(ctx context.Context, req core.DeleteOrderRequest) error {
	err := o.Store.DeleteOrder(ctx, req)
	if err != nil {
		errResp := errors.ErrResponse{
			Code:      errors.ErrorBadRequest,
			MessageID: "",
			Message:   "Failed to delete order",
			Reason:    err.Error(),
		}
		return errors.Errorf(errResp)
	}

	return nil
}
