package product

import (
	"mime/multipart"
	"net/http"

	"github.com/chizidotdev/shop/api/httpUtil"
	"github.com/chizidotdev/shop/api/middleware"
	"github.com/chizidotdev/shop/repository"
	"github.com/gin-gonic/gin"
)

type updateProductRequest struct {
	Title       *string                 `form:"title,omitempty"`
	Description *string                 `form:"description,omitempty"`
	Price       *float64                `form:"price,omitempty"`
	OutOfStock  *bool                   `form:"outOfStock,omitempty"`
	Images      []*multipart.FileHeader `form:"images,omitempty"`
}

func (p *ProductHandler) UpdateProduct(ctx *gin.Context) {
	var req updateProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		httpUtil.Error(ctx, &httpUtil.ErrorResponse{
			Code:      http.StatusBadRequest,
			MessageID: "",
			Message:   "Invalid product request",
			Reason:    err.Error(),
		})
		return
	}

	user := middleware.GetAuthenticatedUser(ctx)

	productID, err := repository.ParseUUID(ctx.Param(productIDParam))
	if err != nil {
		httpUtil.Error(ctx, &httpUtil.ErrorResponse{
			Code:      http.StatusBadRequest,
			MessageID: "",
			Message:   "Invalid product id",
			Reason:    err.Error(),
		})
		return
	}

	existingProduct, err := p.pgStore.GetProduct(ctx, productID)
	if err != nil {
		httpUtil.Error(ctx, &httpUtil.ErrorResponse{
			Code:      http.StatusNotFound,
			MessageID: "",
			Message:   "Product not found",
			Reason:    err.Error(),
		})
		return
	}

	if req.Title != nil && *req.Title != existingProduct.Title {
		existingProduct.Title = *req.Title
	}
	if req.Description != nil && *req.Description != existingProduct.Description {
		existingProduct.Description = *req.Description
	}
	if req.Price != nil && *req.Price != existingProduct.Price {
		existingProduct.Price = *req.Price
	}
	if req.OutOfStock != nil {
		existingProduct.OutOfStock = *req.OutOfStock
	}

	product, err := p.pgStore.UpdateProduct(ctx, repository.UpdateProductParams{
		ID:          productID,
		StoreID:     user.StoreID.UUID,
		Title:       existingProduct.Title,
		Description: existingProduct.Description,
		Price:       existingProduct.Price,
		OutOfStock:  existingProduct.OutOfStock,
	})
	if err != nil {
		httpUtil.Error(ctx, &httpUtil.ErrorResponse{
			Code:      http.StatusInternalServerError,
			MessageID: "",
			Message:   "Failed to update product",
			Reason:    err.Error(),
		})
		return
	}

	httpUtil.Success(ctx, &httpUtil.SuccessResponse{
		Code:    http.StatusOK,
		Message: "Product updated successfully",
		Data:    product,
	})
}
