package product

import (
	"net/http"

	"github.com/chizidotdev/shop/api/httpUtil"
	"github.com/chizidotdev/shop/repository"
	"github.com/gin-gonic/gin"
)

func (p *ProductHandler) ListUserProducts(ctx *gin.Context) {
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

	products, err := p.pgStore.ListProductsByStore(ctx, storeId)
	if err != nil {
		httpUtil.Error(ctx, &httpUtil.ErrorResponse{
			Code:      http.StatusInternalServerError,
			MessageID: "",
			Message:   "Failed to get user products",
			Reason:    err.Error(),
		})
		return
	}

	productResponse := make([]createProductResponse, len(products))
	for i, product := range products {
		productResponse[i].Product = &product
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
		productResponse[i].Images = images
	}

	httpUtil.Success(ctx, &httpUtil.SuccessResponse{
		Code:    http.StatusOK,
		Data:    productResponse,
		Message: "User products retrieved successfully",
	})
}
