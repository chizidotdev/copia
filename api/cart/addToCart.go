package cart

import (
	"net/http"

	"github.com/chizidotdev/shop/api/httpUtil"
	"github.com/chizidotdev/shop/api/middleware"
	"github.com/chizidotdev/shop/repository"
	"github.com/gin-gonic/gin"
)

type addToCartReq struct {
	ProductID string `json:"productId" binding:"required"`
	Quantity  int32  `json:"quantity" binding:"required,gt=0"`
}

func (c *CartHandler) AddToCart(ctx *gin.Context) {
	user := middleware.GetAuthenticatedUser(ctx)

	var req addToCartReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		httpUtil.Error(ctx, &httpUtil.ErrorResponse{
			Code:      http.StatusBadRequest,
			Message:   "Invalid request",
			MessageID: "",
			Reason:    err.Error(),
		})
		return
	}

	productID, err := repository.ParseUUID(req.ProductID)
	if err != nil {
		httpUtil.Error(ctx, &httpUtil.ErrorResponse{
			Code:      http.StatusBadRequest,
			Message:   "Invalid product ID",
			MessageID: "",
			Reason:    err.Error(),
		})
		return
	}

	cartItem, err := c.pgStore.UpsertCartItem(ctx, repository.UpsertCartItemParams{
		UserID:    user.ID,
		ProductID: productID,
		Quantity:  req.Quantity,
	})
	if err != nil {
		httpUtil.Error(ctx, &httpUtil.ErrorResponse{
			Code:      http.StatusInternalServerError,
			Message:   "failed to add item to cart",
			MessageID: "",
			Reason:    err.Error(),
		})
		return
	}

	httpUtil.Success(ctx, &httpUtil.SuccessResponse{
		Data:    cartItem,
		Message: "Item added to cart successfully",
		Code:    http.StatusCreated,
	})
}
