package product

import (
	"net/http"

	"github.com/chizidotdev/shop/api/httpUtil"
	"github.com/chizidotdev/shop/api/middleware"
	"github.com/chizidotdev/shop/repository"
	"github.com/chizidotdev/shop/repository/adapters"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	storeIDParam   = "storeID"
	productIDParam = "productID"
)

type ProductHandler struct {
	pgStore *repository.Repository
	s3Store *adapters.S3Store
}

func NewProductHandler(pgStore *repository.Repository) *ProductHandler {
	return &ProductHandler{pgStore: pgStore}
}

func (p *ProductHandler) validateStorePermissions(ctx *gin.Context) uuid.UUID {
	storeID, err := repository.ParseUUID(ctx.Param(storeIDParam))
	if err != nil {
		httpUtil.Error(ctx, &httpUtil.ErrorResponse{
			Code:      http.StatusBadRequest,
			MessageID: "",
			Message:   "Invalid store id",
			Reason:    err.Error(),
		})
		return uuid.Nil
	}

	store, err := p.pgStore.GetStore(ctx, storeID)
	if err != nil {
		httpUtil.Error(ctx, &httpUtil.ErrorResponse{
			Code:      http.StatusNotFound,
			MessageID: "",
			Message:   "Store not found",
			Reason:    err.Error(),
		})
		return uuid.Nil
	}

	user := middleware.GetAuthenticatedUser(ctx)
	if store.UserID != user.ID {
		httpUtil.Error(ctx, &httpUtil.ErrorResponse{
			Code:      http.StatusForbidden,
			MessageID: "",
			Message:   "You are not authorized to create product for this store",
			Reason:    "",
		})
		return uuid.Nil
	}

	return storeID
}
