package store

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/chizidotdev/shop/api/httpUtil"
	"github.com/chizidotdev/shop/api/middleware"
	"github.com/chizidotdev/shop/repository"
	"github.com/gin-gonic/gin"
)

func (u *StoreHandler) GetStore(ctx *gin.Context) {
	storeId, err := repository.ParseUUID(ctx.Param(storeIDParam))
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

func (u *StoreHandler) GetUserStore(ctx *gin.Context) {
	user := middleware.GetAuthenticatedUser(ctx)

	store, err := u.pgStore.GetStoreByUserId(ctx, user.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println(store)
			httpUtil.Error(ctx, &httpUtil.ErrorResponse{
				Code:      http.StatusNotFound,
				MessageID: "",
				Message:   "Store not found",
				Reason:    err.Error(),
			})
			return
		}

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
