package store

import (
	"net/http"

	"github.com/chizidotdev/shop/api/httpUtil"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (u *StoreHandler) GetStore(ctx *gin.Context) {
	storeId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		errResp := httpUtil.HttpError{
			Code:      httpUtil.ErrorBadRequest,
			MessageID: "",
			Message:   "Invalid store id",
			Reason:    err.Error(),
		}
		httpUtil.Error(ctx, httpUtil.Errorf(errResp))
		return
	}

	store, err := u.pgStore.GetStore(ctx, storeId)
	if err != nil {
		errResp := httpUtil.HttpError{
			Code:      httpUtil.ErrorInternal,
			MessageID: "",
			Message:   "Failed to retrieve store",
			Reason:    err.Error(),
		}
		httpUtil.Error(ctx, httpUtil.Errorf(errResp))
		return
	}

	httpUtil.Success(ctx, http.StatusCreated, httpUtil.SuccessResponse{
		Data:    store,
		Message: "Store retrieved successfully",
	})
}
