package store

import (
	"net/http"

	"github.com/chizidotdev/shop/api/httpUtil"
	"github.com/chizidotdev/shop/api/middleware"
	"github.com/chizidotdev/shop/repository"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createStoreRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

func (s *StoreHandler) CreateStore(ctx *gin.Context) {
	var store createStoreRequest
	err := ctx.BindJSON(&store)
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

	user := middleware.GetAuthenticatedUser(ctx)
	storeProfile, err := s.pgStore.CreateStore(ctx, repository.CreateStoreParams{
		Name:        store.Name,
		Description: store.Description,
		UserID:      user.ID,
	})

	if err != nil {
		code := httpUtil.ErrorInternal
		message := "Failed to create store"

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
		Message: "Store created successfully",
	})
}
