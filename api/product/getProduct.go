package product

import (
	"net/http"

	"github.com/chizidotdev/shop/api/httpUtil"
	"github.com/chizidotdev/shop/repository"
	"github.com/gin-gonic/gin"
)

type getProductResponse struct {
	*repository.Product
	Images *[]repository.ProductImage `json:"images"`
	Store  *repository.Store          `json:"store"`
}

func (p *ProductHandler) GetProduct(ctx *gin.Context) {
	productId, err := repository.ParseUUID(ctx.Param(productIDParam))
	if err != nil {
		httpUtil.Error(ctx, &httpUtil.ErrorResponse{
			Code:      http.StatusBadRequest,
			MessageID: "",
			Message:   "Invalid product id",
			Reason:    err.Error(),
		})
		return
	}

	product, err := p.pgStore.GetProduct(ctx, productId)
	if err != nil {
		httpUtil.Error(ctx, &httpUtil.ErrorResponse{
			Code:      http.StatusNotFound,
			MessageID: "",
			Message:   "Failed to retrieve product",
			Reason:    err.Error(),
		})
		return
	}

	store, err := p.pgStore.GetStore(ctx, product.StoreID)
	if err != nil {
		httpUtil.Error(ctx, &httpUtil.ErrorResponse{
			Code:      http.StatusNotFound,
			MessageID: "",
			Message:   "Failed to retrieve store",
			Reason:    err.Error(),
		})
		return
	}

	images, err := p.pgStore.ListProductImages(ctx, product.ID)
	if err != nil {
		httpUtil.Error(ctx, &httpUtil.ErrorResponse{
			Code:      http.StatusInternalServerError,
			MessageID: "",
			Message:   "Failed to get product images",
			Reason:    err.Error(),
		})
		return
	}

	productResponse := &getProductResponse{
		Product: &product,
		Store:   &store,
		Images:  &images,
	}

	httpUtil.Success(ctx, &httpUtil.SuccessResponse{
		Code:    http.StatusOK,
		Data:    productResponse,
		Message: "Product retrieved successfully",
	})
}
