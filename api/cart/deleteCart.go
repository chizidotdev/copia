package cart

import (
	"net/http"

	"github.com/chizidotdev/shop/api/httpUtil"
	"github.com/chizidotdev/shop/api/middleware"
	"github.com/chizidotdev/shop/repository"
	"github.com/gin-gonic/gin"
)

func (c *CartHandler) DeleteCart(ctx *gin.Context) {
	cartID, err := repository.ParseUUID(ctx.Param(cartIDParam))
	if err != nil {
		httpUtil.Error(ctx, &httpUtil.ErrorResponse{
			Code:      http.StatusBadRequest,
			Message:   "Invalid cart ID",
			MessageID: "",
			Reason:    err.Error(),
		})
		return
	}

	user := middleware.GetAuthenticatedUser(ctx)
	cart, err := c.pgStore.DeleteCartItem(ctx, repository.DeleteCartItemParams{
		UserID: user.ID,
		ID:     cartID,
	})
	if err != nil {
		httpUtil.Error(ctx, &httpUtil.ErrorResponse{
			Code:      http.StatusInternalServerError,
			Message:   "failed to delete item from cart",
			MessageID: "",
			Reason:    err.Error(),
		})
		return
	}

	httpUtil.Success(ctx, &httpUtil.SuccessResponse{
		Data:    cart,
		Message: "Item deleted from cart successfully",
		Code:    http.StatusOK,
	})
}
