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
		errResp := httpUtil.HttpError{
			Code:      httpUtil.ErrorInternal,
			MessageID: "",
			Message:   "Error retrieving stores",
			Reason:    err.Error(),
		}
		httpUtil.Error(ctx, httpUtil.Errorf(errResp))
		return
	}

	httpUtil.Success(ctx, http.StatusCreated, httpUtil.SuccessResponse{
		Data:    store,
		Message: "Stores retrieved successfully",
	})
}
