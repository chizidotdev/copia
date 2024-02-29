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
		httpUtil.Error(ctx, &httpUtil.ErrorResponse{
			Code:      http.StatusBadRequest,
			MessageID: "",
			Message:   "Invalid store id",
			Reason:    err.Error(),
		})
		return
	}

	store, err := u.pgStore.GetStore(ctx, storeId)
	if err != nil {
		httpUtil.Error(ctx, &httpUtil.ErrorResponse{
			Code:      http.StatusInternalServerError,
			MessageID: "",
			Message:   "Failed to retrieve store",
			Reason:    err.Error(),
		})
		return
	}

	httpUtil.Success(ctx, &httpUtil.SuccessResponse{
		Code:    http.StatusOK,
		Data:    store,
		Message: "Store retrieved successfully",
	})
}
