package product

import (
	"net/http"

	"github.com/chizidotdev/shop/api/httpUtil"
	"github.com/chizidotdev/shop/api/middleware"
	"github.com/chizidotdev/shop/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (p *ProductHandler) ListStoreProducts(ctx *gin.Context) {
	user := middleware.GetAuthenticatedUser(ctx)

	var storeID uuid.UUID
	if user.Role == repository.UserRoleVendor {
		storeID = user.StoreID.UUID
	} else {
		var err error
		storeID, err = repository.ParseUUID(ctx.Param(storeIDParam))
		if err != nil {
			httpUtil.Error(ctx, &httpUtil.ErrorResponse{
				Code:      http.StatusBadRequest,
				MessageID: "",
				Message:   "Invalid store id",
				Reason:    err.Error(),
			})
			return
		}
	}

	products, err := p.pgStore.ListProductsByStore(ctx, storeID)
	if err != nil {
		httpUtil.Error(ctx, &httpUtil.ErrorResponse{
			Code:      http.StatusInternalServerError,
			MessageID: "",
			Message:   "Failed to get products",
			Reason:    err.Error(),
		})
		return
	}

	productResponse := make([]createProductResponse, len(products))
	err = p.pgStore.ExecTx(ctx, func(tx *repository.Queries) error {
		for i, product := range products {
			productResponse[i].Product = &product
			images, err := p.pgStore.ListProductImages(ctx, product.ID)
			if err != nil {
				return err
			}
			productResponse[i].Images = images
		}

		return nil
	})
	if err != nil {
		httpUtil.Error(ctx, &httpUtil.ErrorResponse{
			Code:      http.StatusInternalServerError,
			MessageID: "",
			Message:   "Failed to get product images",
			Reason:    err.Error(),
		})
		return
	}

	httpUtil.Success(ctx, &httpUtil.SuccessResponse{
		Code:    http.StatusOK,
		Data:    productResponse,
		Message: "User products retrieved successfully",
	})
}
