package service

import (
	"context"
	"github.com/chizidotdev/copia/dto"
	"github.com/chizidotdev/copia/repository"

	"github.com/google/uuid"
)

type OrderService interface {
	ListOrders(ctx context.Context, userID string) ([]repository.Order, error)
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
		order, err = store.CreateOrder(ctx, repository.CreateOrderParams{
			UserEmail:             req.UserEmail,
			CustomerID:            req.CustomerID,
			Status:                req.Status,
			ShippingDetails:       req.ShippingDetails,
			EstimatedDeliveryDate: req.EstimatedDeliveryDate,
			OrderDate:             req.OrderDate,
			TotalAmount:           req.TotalAmount,
			PaymentStatus:         req.PaymentStatus,
			PaymentMethod:         req.PaymentMethod,
			BillingAddress:        req.BillingAddress,
			ShippingAddress:       req.ShippingAddress,
			Notes:                 req.Notes,
		})
		if err != nil {
			return err
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
				return err
			}
		}

		return err
	})
	if txErr != nil {
		return repository.Order{}, txErr
	}

	return order, nil
}

func (o *orderService) ListOrders(ctx context.Context, userEmail string) ([]repository.Order, error) {
	orders, err := o.Store.ListOrders(ctx, userEmail)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (o *orderService) UpdateOrder(ctx context.Context, req dto.Order) (repository.Order, error) {
	order, err := o.Store.UpdateOrder(ctx, repository.UpdateOrderParams{
		ID:                    req.ID,
		UserEmail:             req.UserEmail,
		CustomerID:            req.CustomerID,
		Status:                req.Status,
		ShippingDetails:       req.ShippingDetails,
		EstimatedDeliveryDate: req.EstimatedDeliveryDate,
		OrderDate:             req.OrderDate,
		TotalAmount:           req.TotalAmount,
		PaymentStatus:         req.PaymentStatus,
		PaymentMethod:         req.PaymentMethod,
		BillingAddress:        req.BillingAddress,
		ShippingAddress:       req.ShippingAddress,
		Notes:                 req.Notes,
	})
	if err != nil {
		return repository.Order{}, err
	}

	return order, nil
}

func (o *orderService) GetOrderByID(ctx context.Context, orderID uuid.UUID) (repository.Order, error) {
	order, err := o.Store.GetOrder(ctx, orderID)
	if err != nil {
		return repository.Order{}, err
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
