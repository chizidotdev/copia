package service

import (
	"context"
	"github.com/chizidotdev/copia/dto"
	"github.com/chizidotdev/copia/repository"
	"github.com/chizidotdev/copia/util"

	"github.com/google/uuid"
)

type OrderService interface {
	ListOrders(ctx context.Context, userID uuid.UUID) ([]repository.Order, error)
	CreateOrder(ctx context.Context, req dto.Order) (repository.Order, error)
	UpdateOrder(ctx context.Context, req dto.Order) (repository.Order, error)
	GetOrderByID(ctx context.Context, orderID uuid.UUID) (repository.Order, error)
	DeleteOrder(ctx context.Context, req dto.DeleteOrderParams) error
}

type orderService struct {
	Store *repository.Repository
}

func NewOrderService(store *repository.Repository) OrderService {
	return &orderService{
		Store: store,
	}
}

func (o *orderService) CreateOrder(ctx context.Context, req dto.Order) (repository.Order, error) {
	var order repository.Order
	txErr := o.Store.ExecTx(ctx, func(store *repository.Repository) error {
		var err error
		totalAmount := float32(0)
		for _, orderItem := range req.OrderItems {
			totalAmount += orderItem.SubTotal
		}
		order, err = store.CreateOrder(ctx, repository.CreateOrderParams{
			UserID:                req.UserID,
			CustomerID:            req.CustomerID,
			Status:                req.Status,
			EstimatedDeliveryDate: req.EstimatedDeliveryDate,
			OrderDate:             req.OrderDate,
			TotalAmount:           totalAmount,
		})
		if err != nil {
			return util.Errorf(util.ErrorBadRequest, "Failed to create order")
		}

		for _, orderItem := range req.OrderItems {
			orderItem.OrderID = order.ID
			_, err = store.CreateOrderItem(ctx, repository.CreateOrderItemParams{
				OrderID:   orderItem.OrderID,
				ProductID: orderItem.ProductID,
				Quantity:  orderItem.Quantity,
				UnitPrice: orderItem.UnitPrice,
				SubTotal:  orderItem.SubTotal,
			})
			if err != nil {
				return util.Errorf(util.ErrorBadRequest, "Failed to create order item")
			}
		}

		return err
	})
	if txErr != nil {
		return repository.Order{}, txErr
	}

	return order, nil
}

func (o *orderService) ListOrders(ctx context.Context, userID uuid.UUID) ([]repository.Order, error) {
	orders, err := o.Store.ListOrders(ctx, userID)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (o *orderService) UpdateOrder(ctx context.Context, req dto.Order) (repository.Order, error) {
	var order repository.Order
	txErr := o.Store.ExecTx(ctx, func(store *repository.Repository) error {
		totalAmount := float32(0)
		for _, orderItem := range req.OrderItems {
			totalAmount += orderItem.SubTotal
		}
		var err error
		order, err = o.Store.UpdateOrder(ctx, repository.UpdateOrderParams{
			ID:                    req.ID,
			UserID:                req.UserID,
			CustomerID:            req.CustomerID,
			Status:                req.Status,
			EstimatedDeliveryDate: req.EstimatedDeliveryDate,
			OrderDate:             req.OrderDate,
			TotalAmount:           totalAmount,
		})
		if err != nil {
			return util.Errorf(util.ErrorNotFound, "Failed to update order")
		}

		for _, orderItem := range req.OrderItems {
			orderItem.OrderID = order.ID
			_, err = store.UpdateOrderItem(ctx, repository.UpdateOrderItemParams{
				ID:        orderItem.ID,
				OrderID:   orderItem.OrderID,
				ProductID: orderItem.ProductID,
				Quantity:  orderItem.Quantity,
				UnitPrice: orderItem.UnitPrice,
				SubTotal:  orderItem.SubTotal,
			})
			if err != nil {
				return util.Errorf(util.ErrorNotFound, "Failed to update order item")
			}
		}

		return err
	})
	if txErr != nil {
		return repository.Order{}, txErr
	}

	return order, nil
}

func (o *orderService) GetOrderByID(ctx context.Context, orderID uuid.UUID) (repository.Order, error) {
	order, err := o.Store.GetOrder(ctx, orderID)
	if err != nil {
		return repository.Order{}, util.Errorf(util.ErrorNotFound, "Order not found")
	}

	return order, nil
}

func (o *orderService) DeleteOrder(ctx context.Context, req dto.DeleteOrderParams) error {
	err := o.Store.DeleteOrder(ctx, repository.DeleteOrderParams(req))
	if err != nil {
		return err
	}

	return nil
}

// TODO: Implement
// func (o *orderService) DeleteOrderItem
