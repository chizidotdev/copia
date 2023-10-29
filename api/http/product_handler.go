package http

import (
	"github.com/chizidotdev/copia/api/http/middleware"
	"github.com/chizidotdev/copia/internal/app/core"
	"github.com/chizidotdev/copia/internal/app/usecases"
	"github.com/chizidotdev/copia/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type ProductHandler struct {
	ProductService *usecases.ProductService
}

func NewProductHandler(productService *usecases.ProductService) *ProductHandler {
	return &ProductHandler{
		ProductService: productService,
	}
}

func (p *ProductHandler) createProduct(ctx *gin.Context) {
	var req core.ProductRequest
	if err := ctx.ShouldBind(&req); err != nil {
		errorResponse(ctx, errors.Errorf(errors.ErrorBadRequest, "Invalid request payload."))
		return
	}

	user := middleware.GetAuthenticatedUser(ctx)
	product, err := p.ProductService.CreateProduct(ctx, core.ProductRequest{
		UserID:          user.ID,
		Name:            req.Name,
		Description:     req.Description,
		Price:           req.Price,
		QuantityInStock: req.QuantityInStock,
		Image:           req.Image,
	})
	if err != nil {
		errorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, product)
}

func (p *ProductHandler) updateProduct(ctx *gin.Context) {
	IDParam := ctx.Param("id")
	productID, err := uuid.Parse(IDParam)
	if err != nil {
		errorResponse(ctx, errors.Errorf(errors.ErrorBadRequest, "Invalid product ID"))
		return
	}

	var req core.Product
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errorResponse(ctx, errors.Errorf(errors.ErrorBadRequest, "Invalid request payload."))
		return
	}

	user := middleware.GetAuthenticatedUser(ctx)
	req.UserID = user.ID
	product, err := p.ProductService.UpdateProduct(ctx, core.Product{
		ID:              productID,
		UserID:          user.ID,
		Name:            req.Name,
		Description:     req.Description,
		Price:           req.Price,
		QuantityInStock: req.QuantityInStock,
		ImageURL:        req.ImageURL,
	})
	if err != nil {
		errorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, product)
}

func (p *ProductHandler) updateProductImage(ctx *gin.Context) {
	IDParam := ctx.Param("id")
	productID, err := uuid.Parse(IDParam)
	if err != nil {
		errorResponse(ctx, errors.Errorf(errors.ErrorBadRequest, "Invalid product ID"))
		return
	}

	var req core.ProductImageRequest
	if err := ctx.ShouldBind(&req); err != nil {
		errorResponse(ctx, errors.Errorf(errors.ErrorBadRequest, "Invalid request payload."))
		return
	}

	user := middleware.GetAuthenticatedUser(ctx)
	req.UserID = user.ID
	product, err := p.ProductService.UpdateProductImage(ctx, core.ProductImageRequest{
		ID:     productID,
		UserID: user.ID,
		Image:  req.Image,
	})
	if err != nil {
		errorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, product)
}

func (p *ProductHandler) listProducts(ctx *gin.Context) {
	user := middleware.GetAuthenticatedUser(ctx)
	products, err := p.ProductService.ListProducts(ctx, user.ID)
	if err != nil {
		errorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, products)
}

func (p *ProductHandler) getProduct(ctx *gin.Context) {
	IDParam := ctx.Param("id")
	productID, err := uuid.Parse(IDParam)
	if err != nil {
		errorResponse(ctx, errors.Errorf(errors.ErrorBadRequest, "Invalid product ID"))
		return
	}

	product, err := p.ProductService.GetProductByID(ctx, productID)
	if err != nil {
		errorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, product)
}

func (p *ProductHandler) deleteProduct(ctx *gin.Context) {
	IDParam := ctx.Param("id")
	productID, err := uuid.Parse(IDParam)
	if err != nil {
		errorResponse(ctx, errors.Errorf(errors.ErrorBadRequest, "Invalid product ID"))
		return
	}

	user := middleware.GetAuthenticatedUser(ctx)
	err = p.ProductService.DeleteProduct(ctx, core.DeleteProductRequest{
		UserID: user.ID,
		ID:     productID,
	})
	if err != nil {
		errorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, "Successfully deleted product.")
}

func (p *ProductHandler) updateProductQuantity(ctx *gin.Context) {
	IDParam := ctx.Param("id")
	productID, err := uuid.Parse(IDParam)
	if err != nil {
		errorResponse(ctx, errors.Errorf(errors.ErrorBadRequest, "Invalid product ID"))
		return
	}

	var req core.UpdateProductQuantityRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errorResponse(ctx, errors.Errorf(errors.ErrorBadRequest, "Invalid request payload."))
		return
	}

	user := middleware.GetAuthenticatedUser(ctx)
	product, err := p.ProductService.UpdateProductQuantity(ctx, core.UpdateProductQuantityRequest{
		UserID:      user.ID,
		ID:          productID,
		NewQuantity: req.NewQuantity,
	})
	if err != nil {
		errorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, product)
}

func (p *ProductHandler) updateProductSettings(ctx *gin.Context) {
	var req core.ProductSettings
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errorResponse(ctx, errors.Errorf(errors.ErrorBadRequest, "Invalid request payload."))
		return
	}

	user := middleware.GetAuthenticatedUser(ctx)
	settings, err := p.ProductService.UpdateProductSettings(ctx, core.ProductSettings{
		UserID:       user.ID,
		ReorderPoint: req.ReorderPoint,
	})
	if err != nil {
		errorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, settings)
}