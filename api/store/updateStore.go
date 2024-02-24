package store

import (
	"net/http"

	"github.com/chizidotdev/shop/api/httpUtil"
	"github.com/chizidotdev/shop/api/middleware"
	"github.com/chizidotdev/shop/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type updateStoreRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

func (s *StoreHandler) UpdateStore(ctx *gin.Context) {
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

	var store updateStoreRequest
	err = ctx.BindJSON(&store)
	if err != nil {
		errResp := httpUtil.HttpError{
			Code:      httpUtil.ErrorBadRequest,
			MessageID: "",
			Message:   "Invalid store request",
			Reason:    err.Error(),
		}
		httpUtil.Error(ctx, httpUtil.Errorf(errResp))
		return
	}

	existingStore, err := s.pgStore.GetStore(ctx, storeId)
	if err != nil {
		errResp := httpUtil.HttpError{
			Code:      httpUtil.ErrorNotFound,
			MessageID: "",
			Message:   "Store not found",
			Reason:    err.Error(),
		}
		httpUtil.Error(ctx, httpUtil.Errorf(errResp))
		return
	}

	user := middleware.GetAuthenticatedUser(ctx)
	if existingStore.UserID != user.ID {
		errResp := httpUtil.HttpError{
			Code:      httpUtil.ErrorForbidden,
			MessageID: "",
			Message:   "You are not authorized to update this store",
			Reason:    err.Error(),
		}
		httpUtil.Error(ctx, httpUtil.Errorf(errResp))
		return
	}

	storeProfile, err := s.pgStore.UpdateStore(ctx, repository.UpdateStoreParams{
		ID:          storeId,
		Name:        store.Name,
		Description: store.Description,
	})

	if err != nil {
		code := httpUtil.ErrorInternal
		message := "Failed to update store"

		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation", "foreign_key_violation":
				code = httpUtil.ErrorForbidden
				message = "Store name already exists"
			}
		}

		errResp := httpUtil.HttpError{
			Code:      code,
			MessageID: "",
			Message:   message,
			Reason:    err.Error(),
		}
		httpUtil.Error(ctx, httpUtil.Errorf(errResp))
		return
	}

	httpUtil.Success(ctx, http.StatusCreated, httpUtil.SuccessResponse{
		Data:    storeProfile,
		Message: "Store updated successfully",
	})
}
