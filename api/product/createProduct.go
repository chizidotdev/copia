package product

import (
	"mime/multipart"
	"net/http"

	"github.com/chizidotdev/shop/api/httpUtil"
	"github.com/chizidotdev/shop/repository"
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
	Images []repository.ProductImage
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

	storeID := p.validateStorePermissions(ctx)

	var newProduct repository.Product
	err := p.pgStore.ExecTx(ctx, func(tx *repository.Queries) error {
		var txErr error
		newProduct, txErr = p.pgStore.CreateProduct(ctx, repository.CreateProductParams{
			StoreID:     storeID,
			Title:       product.Title,
			Description: product.Description,
			Price:       product.Price,
			OutOfStock:  false,
		})
		if txErr != nil {
			return txErr
		}

		if len(product.Images) > 0 {
			/*image, err := util.ParseImage(product.Images[0])
			if err != nil {
				return err
			}

			imageUrl, err := p.s3Store.UploadFile(newProduct.ID.String(), image)
			if err != nil {
				return err
			}*/

			_, txErr = p.pgStore.CreateProductImage(ctx, repository.CreateProductImageParams{
				ProductID: newProduct.ID,
				Url:       "",
			})
			if txErr != nil {
				return txErr
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
		Images:  []repository.ProductImage{},
	}

	httpUtil.Success(ctx, &httpUtil.SuccessResponse{
		Code:    http.StatusCreated,
		Message: "Product created successfully",
		Data:    resp,
	})
}
