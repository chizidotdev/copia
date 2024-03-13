package order

import (
	"net/http"

	"github.com/chizidotdev/shop/api/httpUtil"
	"github.com/chizidotdev/shop/api/middleware"
	"github.com/chizidotdev/shop/repository"
	"github.com/gin-gonic/gin"
)

func (c *OrderHandler) DeleteOrder(ctx *gin.Context) {
	orderID, err := repository.ParseUUID(ctx.Param(orderIDParam))
	if err != nil {
		httpUtil.Error(ctx, &httpUtil.ErrorResponse{
			Code:      http.StatusBadRequest,
			Message:   "Invalid order ID",
			MessageID: "",
			Reason:    err.Error(),
		})
		return
	}

	user := middleware.GetAuthenticatedUser(ctx)
	order, err := c.pgStore.DeleteOrder(ctx, repository.DeleteOrderParams{
		UserID: user.ID,
		ID:     orderID,
	})
	if err != nil {
		httpUtil.Error(ctx, &httpUtil.ErrorResponse{
			Code:      http.StatusInternalServerError,
			Message:   "failed to delete order",
			MessageID: "",
			Reason:    err.Error(),
		})
		return
	}

	httpUtil.Success(ctx, &httpUtil.SuccessResponse{
		Data:    order,
		Message: "Item deleted from cart successfully",
		Code:    http.StatusOK,
	})
}
