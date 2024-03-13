package store

import (
	"net/http"

	"github.com/chizidotdev/shop/api/httpUtil"
	"github.com/chizidotdev/shop/api/middleware"
	"github.com/chizidotdev/shop/repository"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type updateStoreRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

func (s *StoreHandler) UpdateUserStore(ctx *gin.Context) {
	var store updateStoreRequest
	err := ctx.BindJSON(&store)
	if err != nil {
		httpUtil.Error(ctx, &httpUtil.ErrorResponse{
			Code:      http.StatusBadRequest,
			MessageID: "",
			Message:   "Invalid store request",
			Reason:    err.Error(),
		})
		return
	}

	user := middleware.GetAuthenticatedUser(ctx)

	storeProfile, err := s.pgStore.UpdateStore(ctx, repository.UpdateStoreParams{
		ID:          user.StoreID.UUID,
		Name:        store.Name,
		Description: store.Description,
	})
	if err != nil {
		code := http.StatusInternalServerError
		message := "Failed to update store"

		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation", "foreign_key_violation":
				code = http.StatusForbidden
				message = "Store name already exists"
			}
		}

		httpUtil.Error(ctx, &httpUtil.ErrorResponse{
			Code:      code,
			MessageID: "",
			Message:   message,
			Reason:    err.Error(),
		})
		return
	}

	httpUtil.Success(ctx, &httpUtil.SuccessResponse{
		Code:    http.StatusOK,
		Data:    storeProfile,
		Message: "Store updated successfully",
	})
}
