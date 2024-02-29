package store

import (
	"net/http"

	"github.com/chizidotdev/shop/api/httpUtil"
	"github.com/chizidotdev/shop/api/middleware"
	"github.com/gin-gonic/gin"
)

func (u *StoreHandler) ListStores(ctx *gin.Context) {
	user := middleware.GetAuthenticatedUser(ctx)

	store, err := u.pgStore.ListStores(ctx, user.ID)
	if err != nil {
		httpUtil.Error(ctx, &httpUtil.ErrorResponse{
			Code:      http.StatusInternalServerError,
			MessageID: "",
			Message:   "Error retrieving stores",
			Reason:    err.Error(),
		})
		return
	}

	httpUtil.Success(ctx, &httpUtil.SuccessResponse{
		Code:    http.StatusOK,
		Data:    store,
		Message: "Stores retrieved successfully",
	})
}
