package product

import (
	"net/http"

	"github.com/chizidotdev/shop/api/httpUtil"
	"github.com/chizidotdev/shop/repository"
	"github.com/gin-gonic/gin"
)

func (p *ProductHandler) DeleteProduct(ctx *gin.Context) {
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

	err = p.pgStore.DeleteProduct(ctx, repository.DeleteProductParams{
		ID:      productID,
		StoreID: storeID,
	})
	if err != nil {
		httpUtil.Error(ctx, &httpUtil.ErrorResponse{
			Code:      http.StatusInternalServerError,
			MessageID: "",
			Message:   "Failed to delete product",
			Reason:    err.Error(),
		})
		return
	}

	httpUtil.Success(ctx, &httpUtil.SuccessResponse{
		Code:    http.StatusOK,
		Message: "Product deleted successfully",
	})
}
