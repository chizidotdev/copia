package cart

import (
	"net/http"

	"github.com/chizidotdev/shop/api/httpUtil"
	"github.com/chizidotdev/shop/api/middleware"
	"github.com/chizidotdev/shop/repository"
	"github.com/gin-gonic/gin"
)

type updateCartReq struct {
	Quantity int32 `json:"quantity" binding:"required,gt=0"`
}

func (c *CartHandler) UpdateCart(ctx *gin.Context) {
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

	var req updateCartReq
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
	cartItem, err := c.pgStore.UpdateCartItemQuantity(ctx, repository.UpdateCartItemQuantityParams{
		UserID:   user.ID,
		ID:       cartID,
		Quantity: req.Quantity,
	})
	if err != nil {
		httpUtil.Error(ctx, &httpUtil.ErrorResponse{
			Code:      http.StatusInternalServerError,
			Message:   "failed to update item in cart",
			MessageID: "",
			Reason:    err.Error(),
		})
		return
	}

	httpUtil.Success(ctx, &httpUtil.SuccessResponse{
		Data:    cartItem,
		Message: "Item updated in cart successfully",
		Code:    http.StatusOK,
	})
}
