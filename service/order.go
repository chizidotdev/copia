package service

import (
	"context"
	"github.com/chizidotdev/copia/dto"
	repository2 "github.com/chizidotdev/copia/repository"

	"github.com/google/uuid"
)

type OrderService interface {
	ListOrders(ctx context.Context, userID string) ([]repository2.Order, error)
	CreateOrder(ctx context.Context, req dto.Order) (repository2.Order, error)
	UpdateOrder(ctx context.Context, req dto.Order) (repository2.Order, error)
	GetOrderByID(ctx context.Context, orderID uuid.UUID) (repository2.Order, error)
	DeleteOrder(ctx context.Context, req dto.DeleteOrderParams) error
}

type orderService struct {
	Store *repository2.Store
}

func NewOrderService(store *repository2.Store) OrderService {
	return &orderService{
		Store: store,
	}
}

func (o *orderService) CreateOrder(ctx context.Context, req dto.Order) (repository2.Order, error) {
	var order repository2.Order
	txErr := o.Store.ExecTx(ctx, func(store *repository2.Store) error {
		var err error
		order, err = store.CreateOrder(ctx, repository2.CreateOrderParams{
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
			_, err = store.CreateOrderItem(ctx, repository2.CreateOrderItemParams{
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
		return repository2.Order{}, txErr
	}

	return order, nil
}

func (o *orderService) ListOrders(ctx context.Context, userEmail string) ([]repository2.Order, error) {
	orders, err := o.Store.ListOrders(ctx, userEmail)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (o *orderService) UpdateOrder(ctx context.Context, req dto.Order) (repository2.Order, error) {
	order, err := o.Store.UpdateOrder(ctx, repository2.UpdateOrderParams{
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
		return repository2.Order{}, err
	}

	return order, nil
}

func (o *orderService) GetOrderByID(ctx context.Context, orderID uuid.UUID) (repository2.Order, error) {
	order, err := o.Store.GetOrder(ctx, orderID)
	if err != nil {
		return repository2.Order{}, err
	}

	return order, nil
}

func (o *orderService) DeleteOrder(ctx context.Context, req dto.DeleteOrderParams) error {
	err := o.Store.DeleteOrder(ctx, repository2.DeleteOrderParams(req))
	if err != nil {
		return err
	}

	return nil
}

// TODO: Implement
// func (o *orderService) DeleteOrderItem
