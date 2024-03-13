package order

import (
	"fmt"
	"net/http"

	"github.com/chizidotdev/shop/api/httpUtil"
	"github.com/chizidotdev/shop/api/middleware"
	"github.com/chizidotdev/shop/repository"
	"github.com/gin-gonic/gin"
)

func (c *OrderHandler) ListStoreOrders(ctx *gin.Context) {
	user := middleware.GetAuthenticatedUser(ctx)

	orders, err := c.pgStore.ListStoreOrderItems(ctx, user.StoreID.UUID)
	if err != nil {
		httpUtil.Error(ctx, &httpUtil.ErrorResponse{
			Code:      http.StatusNotFound,
			Message:   "failed to get orders",
			MessageID: "",
			Reason:    err.Error(),
		})
		return
	}

	httpUtil.Success(ctx, &httpUtil.SuccessResponse{
		Data:    orders,
		Message: "Orders retrieved successfully",
		Code:    http.StatusOK,
	})
}

type listOrder struct {
	*repository.Order
	OrderItems *[]repository.ListOrderItemsRow `json:"orderItems"`
}

func (c *OrderHandler) ListUserOrders(ctx *gin.Context) {
	user := middleware.GetAuthenticatedUser(ctx)

	var resp []listOrder
	err := c.pgStore.ExecTx(ctx, func(tx *repository.Queries) error {
		var txErr error
		orders, txErr := c.pgStore.ListUserOrders(ctx, user.ID)
		if txErr != nil {
			return fmt.Errorf("failed to get orders: %w", txErr)
		}

		for _, order := range orders {
			items, txErr := c.pgStore.ListOrderItems(ctx, order.ID)
			if txErr != nil {
				return fmt.Errorf("failed to get order items: %w", txErr)
			}
			o := listOrder{
				Order:      &order,
				OrderItems: &items,
			}

			resp = append(resp, o)
		}

		return nil
	})
	if err != nil {
		httpUtil.Error(ctx, &httpUtil.ErrorResponse{
			Code:      http.StatusNotFound,
			Message:   "failed to get orders",
			MessageID: "",
			Reason:    err.Error(),
		})
	}

	httpUtil.Success(ctx, &httpUtil.SuccessResponse{
		Data:    resp,
		Message: "Orders retrieved successfully",
		Code:    http.StatusOK,
	})

}
