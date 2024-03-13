package store

import (
	"net/http"

	"github.com/chizidotdev/shop/api/httpUtil"
	"github.com/chizidotdev/shop/api/middleware"
	"github.com/chizidotdev/shop/repository"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type createStoreRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

func (s *StoreHandler) CreateStore(ctx *gin.Context) {
	var req createStoreRequest
	err := ctx.BindJSON(&req)
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
	store, err := s.pgStore.CreateStore(ctx, repository.CreateStoreParams{
		Name:        req.Name,
		Description: req.Description,
		UserID:      user.ID,
	})

	if err != nil {
		code := http.StatusInternalServerError
		message := "Failed to create store"

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

	user.StoreID = uuid.NullUUID{UUID: store.ID, Valid: true}
	session := sessions.Default(ctx)
	session.Set(middleware.ProfileKey, user)
	if err := session.Save(); err != nil {
		httpUtil.Error(ctx, &httpUtil.ErrorResponse{
			Code:      http.StatusInternalServerError,
			MessageID: "",
			Message:   "Failed to create store",
			Reason:    err.Error(),
		})
		return
	}

	httpUtil.Success(ctx, &httpUtil.SuccessResponse{
		Code:    http.StatusCreated,
		Data:    store,
		Message: "Store created successfully",
	})
}
