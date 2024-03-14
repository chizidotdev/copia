package order

import (
	"fmt"
	"net/http"
	"time"

	"github.com/chizidotdev/shop/api/httpUtil"
	"github.com/chizidotdev/shop/api/middleware"
	"github.com/chizidotdev/shop/repository"
	"github.com/gin-gonic/gin"
)

type createOrderResponse struct {
	*repository.Order
	OrderItems *[]repository.OrderItem `json:"orderItems"`
}

func (c *OrderHandler) CreateOrder(ctx *gin.Context) {
	user := middleware.GetAuthenticatedUser(ctx)

	cartItems, err := c.pgStore.GetCartItems(ctx, user.ID)
	if err != nil {
		httpUtil.Error(ctx, &httpUtil.ErrorResponse{
			Code:      http.StatusInternalServerError,
			Message:   "failed to get user cart items",
			MessageID: "",
			Reason:    err.Error(),
		})
		return
	}

	var orderResp createOrderResponse
	err = c.pgStore.ExecTx(ctx, func(tx *repository.Queries) error {
		cartTotal := 0.0
		for _, item := range cartItems {
			cartTotal += item.Price * float64(item.Quantity)
		}

		var txErr error
		order, txErr := c.pgStore.CreateOrder(ctx, repository.CreateOrderParams{
			UserID:          user.ID,
			OrderDate:       time.Now(),
			TotalAmount:     cartTotal,
			PaymentStatus:   repository.PaymentStatusPaid,
			ShippingAddress: "",
		})
		if txErr != nil {
			return fmt.Errorf("failed to create order: %w", txErr)
		}

		var orderItems []repository.OrderItem
		for _, item := range cartItems {
			orderItem, txErr := tx.CreateOrderItem(ctx, repository.CreateOrderItemParams{
				OrderID:   order.ID,
				StoreID:   item.StoreID,
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				UnitPrice: item.Price,
				Subtotal:  item.Price * float64(item.Quantity),
				Status:    repository.OrderStatusPending,
			})
			if txErr != nil {
				return fmt.Errorf("failed to create order item: %w", txErr)
			}

			orderItems = append(orderItems, orderItem)
		}

		orderResp = createOrderResponse{
			Order:      &order,
			OrderItems: &orderItems,
		}

		return nil
	})
	if err != nil {
		httpUtil.Error(ctx, &httpUtil.ErrorResponse{
			Code:      http.StatusInternalServerError,
			Message:   "failed to create order",
			MessageID: "",
			Reason:    err.Error(),
		})
		return
	}

	_ = c.pgStore.ClearCartItems(ctx, user.ID)

	httpUtil.Success(ctx, &httpUtil.SuccessResponse{
		Data:    orderResp,
		Message: "Order created successfully",
		Code:    http.StatusCreated,
	})
}
