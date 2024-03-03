package product

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/chizidotdev/shop/api/httpUtil"
	"github.com/chizidotdev/shop/repository"
	"github.com/gin-gonic/gin"
)

type updateProductRequest struct {
	Title       *string  `json:"title,omitempty"`
	Description *string  `json:"description,omitempty"`
	Price       *float64 `json:"price,omitempty"`
	OutOfStock  *bool    `json:"outOfStock,omitempty"`
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

	storeID := p.validateStorePermissions(ctx)

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

	reqJSON, _ := json.Marshal(req)
	fmt.Println("reqJSON: ", string(reqJSON))
	fmt.Println("req: ", req)

	if req.Title != nil {
		existingProduct.Title = *req.Title
	}
	if req.Description != nil {
		existingProduct.Description = *req.Description
	}
	if req.Price != nil {
		existingProduct.Price = *req.Price
	}
	if req.OutOfStock != nil {
		existingProduct.OutOfStock = *req.OutOfStock
	}

	product, err := p.pgStore.UpdateProduct(ctx, repository.UpdateProductParams{
		ID:          productID,
		StoreID:     storeID,
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
