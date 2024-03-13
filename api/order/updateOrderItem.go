package order

import (
	"net/http"

	"github.com/chizidotdev/shop/api/httpUtil"
	"github.com/chizidotdev/shop/api/middleware"
	"github.com/chizidotdev/shop/repository"
	"github.com/gin-gonic/gin"
)

type updateOrderItemReq struct {
	Status          repository.OrderStatus   `json:"status" binding:"required"`
	PaymentStatus   repository.PaymentStatus `json:"paymentStatus" binding:"required"`
	ShippingAddress string                   `json:"shippingAddress" binding:"required"`
}

func (c *OrderHandler) UpdateOrderItem(ctx *gin.Context) {
	orderID, err := repository.ParseUUID(ctx.Param(orderItemIDParam))
	if err != nil {
		httpUtil.Error(ctx, &httpUtil.ErrorResponse{
			Code:      http.StatusBadRequest,
			Message:   "Invalid order item ID",
			MessageID: "",
			Reason:    err.Error(),
		})
		return
	}

	var req updateOrderItemReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		httpUtil.Error(ctx, &httpUtil.ErrorResponse{
			Code:      http.StatusBadRequest,
			Message:   "Invalid request",
			MessageID: "",
			Reason:    err.Error(),
		})
		return
	}

	user := middleware.GetAuthenticatedUser(ctx)

	orderItem, err := c.pgStore.UpdateOrderItem(ctx, repository.UpdateOrderItemParams{
		ID:      orderID,
		StoreID: user.StoreID.UUID,
		Status:  req.Status,
	})
	if err != nil {
		httpUtil.Error(ctx, &httpUtil.ErrorResponse{
			Code:      http.StatusInternalServerError,
			Message:   "failed to update item status",
			MessageID: "",
			Reason:    err.Error(),
		})
		return
	}

	httpUtil.Success(ctx, &httpUtil.SuccessResponse{
		Data:    orderItem,
		Message: "Order Item updated successfully",
		Code:    http.StatusOK,
	})
}
