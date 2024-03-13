package order

import (
	"net/http"

	"github.com/chizidotdev/shop/api/httpUtil"
	"github.com/chizidotdev/shop/api/middleware"
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

func (c *OrderHandler) ListUserOrders(ctx *gin.Context) {
	user := middleware.GetAuthenticatedUser(ctx)

	orders, err := c.pgStore.ListUserOrders(ctx, user.ID)
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
