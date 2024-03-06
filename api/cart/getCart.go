package cart

import (
	"net/http"

	"github.com/chizidotdev/shop/api/httpUtil"
	"github.com/chizidotdev/shop/api/middleware"
	"github.com/gin-gonic/gin"
)

func (c *CartHandler) GetCart(ctx *gin.Context) {
	user := middleware.GetAuthenticatedUser(ctx)

	cart, err := c.pgStore.GetCartItems(ctx, user.ID)
	if err != nil {
		httpUtil.Error(ctx, &httpUtil.ErrorResponse{
			Code:      http.StatusNotFound,
			Message:   "failed to get cart",
			MessageID: "",
			Reason:    err.Error(),
		})
		return
	}

	httpUtil.Success(ctx, &httpUtil.SuccessResponse{
		Data:    cart,
		Message: "Cart retrieved successfully",
		Code:    http.StatusOK,
	})
}
