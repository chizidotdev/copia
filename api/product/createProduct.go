package product

import (
	"fmt"
	"mime/multipart"
	"net/http"

	"github.com/chizidotdev/shop/api/httpUtil"
	"github.com/chizidotdev/shop/api/middleware"
	"github.com/chizidotdev/shop/repository"
	"github.com/chizidotdev/shop/util"
	"github.com/gin-gonic/gin"
)

type createProductRequest struct {
	Title       string                  `form:"title" binding:"required"`
	Description string                  `form:"description"`
	Price       float64                 `form:"price" binding:"required"`
	OutOfStock  bool                    `form:"outOfStock"`
	Images      []*multipart.FileHeader `form:"images"`
}

type createProductResponse struct {
	*repository.Product
	Images []repository.ProductImage `json:"images"`
}

func (p *ProductHandler) CreateProduct(ctx *gin.Context) {
	var product createProductRequest
	if err := ctx.ShouldBind(&product); err != nil {
		httpUtil.Error(ctx, &httpUtil.ErrorResponse{
			Code:      http.StatusBadRequest,
			MessageID: "",
			Message:   "Invalid product request",
			Reason:    err.Error(),
		})
		return
	}

	user := middleware.GetAuthenticatedUser(ctx)

	var newProduct repository.Product
	var productImages []repository.ProductImage
	err := p.pgStore.ExecTx(ctx, func(tx *repository.Queries) error {
		var txErr error
		newProduct, txErr = p.pgStore.CreateProduct(ctx, repository.CreateProductParams{
			StoreID:     user.StoreID.UUID,
			Title:       product.Title,
			Description: product.Description,
			Price:       product.Price,
			OutOfStock:  false,
		})
		if txErr != nil {
			return txErr
		}

		if len(product.Images) > 0 {
			for _, img := range product.Images {
				image, err := util.ParseImage(img)
				if err != nil {
					return err
				}

				fileName := fmt.Sprintf("%s-%s-%s", user.FirstName, newProduct.ID.String(), image.Name)
				imageUrl, err := p.s3Store.UploadFile(fileName, image.File)
				if err != nil {
					return err
				}

				productImage, txErr := p.pgStore.CreateProductImage(ctx, repository.CreateProductImageParams{
					ProductID: newProduct.ID,
					Url:       imageUrl,
				})
				if txErr != nil {
					return txErr
				}
				productImages = append(productImages, productImage)
			}
		}

		return nil
	})
	if err != nil {
		httpUtil.Error(ctx, &httpUtil.ErrorResponse{
			Code:      http.StatusInternalServerError,
			MessageID: "",
			Message:   "Failed to create product",
			Reason:    err.Error(),
		})
		return
	}

	resp := createProductResponse{
		Product: &newProduct,
		Images:  productImages,
	}

	httpUtil.Success(ctx, &httpUtil.SuccessResponse{
		Code:    http.StatusCreated,
		Message: "Product created successfully",
		Data:    resp,
	})
}
