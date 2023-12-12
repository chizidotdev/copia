package usecases

import (
	"context"
	"github.com/chizidotdev/copia/internal/app/core"
	"github.com/chizidotdev/copia/pkg/errors"
	"github.com/google/uuid"
)

type ProductService struct {
	Store   core.ProductRepository
	s3Store core.FileUploadRepository
}

func NewProductService(productRepo core.ProductRepository, s3Store core.FileUploadRepository) *ProductService {
	return &ProductService{
		Store:   productRepo,
		s3Store: s3Store,
	}
}

func (p *ProductService) CreateProduct(ctx context.Context, req core.ProductRequest) (core.Product, error) {
	sku := core.GenerateSKU(req.UserID.String(), req.Name)

	imageUrl := ""
	if req.Image != nil {
		image, err := core.ParseImage(req.Image)
		if err != nil {
			errResp := errors.ErrResponse{
				Code:      errors.ErrorBadRequest,
				MessageID: "",
				Message:   "Failed to open file",
				Reason:    err.Error(),
			}
			return core.Product{}, errors.Errorf(errResp)
		}

		imageUrl, err = p.s3Store.UploadFile(sku, image)
		if err != nil {
			errResp := errors.ErrResponse{
				Code:      errors.ErrorInternal,
				MessageID: "",
				Message:   "Failed to upload file",
				Reason:    err.Error(),
			}
			return core.Product{}, errors.Errorf(errResp)
		}
	}

	product, err := p.Store.CreateProduct(ctx, core.Product{
		UserID:          req.UserID,
		Name:            req.Name,
		Description:     req.Description,
		Price:           req.Price,
		QuantityInStock: req.QuantityInStock,
		ImageURL:        imageUrl,
		SKU:             sku,
	})
	if err != nil {
		errResp := errors.ErrResponse{
			Code:      errors.ErrorBadRequest,
			MessageID: "",
			Message:   "Failed to create product",
			Reason:    err.Error(),
		}
		return core.Product{}, errors.Errorf(errResp)
	}

	return product, nil
}

func (p *ProductService) UpdateProduct(ctx context.Context, req core.Product) (core.Product, error) {
	product, err := p.Store.UpdateProduct(ctx, req)
	if err != nil {
		errResp := errors.ErrResponse{
			Code:      errors.ErrorBadRequest,
			MessageID: "",
			Message:   "Failed to update product",
			Reason:    err.Error(),
		}
		return core.Product{}, errors.Errorf(errResp)
	}

	return product, nil
}

func (p *ProductService) UpdateProductImage(ctx context.Context, req core.ProductImageRequest) (core.Product, error) {
	product, err := p.Store.GetProduct(ctx, req.ID)
	if err != nil {
		errResp := errors.ErrResponse{
			Code:      errors.ErrorNotFound,
			MessageID: "",
			Message:   "Failed to update product",
			Reason:    err.Error(),
		}
		return core.Product{}, errors.Errorf(errResp)
	}

	image, err := core.ParseImage(req.Image)
	if err != nil {
		errResp := errors.ErrResponse{
			Code:      errors.ErrorBadRequest,
			MessageID: "",
			Message:   "Failed to open file",
			Reason:    err.Error(),
		}
		return core.Product{}, errors.Errorf(errResp)
	}

	imageUrl, err := p.s3Store.UploadFile(product.SKU, image)
	if err != nil {
		errResp := errors.ErrResponse{
			Code:      errors.ErrorInternal,
			MessageID: "",
			Message:   "Failed to upload file",
			Reason:    err.Error(),
		}
		return core.Product{}, errors.Errorf(errResp)
	}

	product.ImageURL = imageUrl

	product, err = p.Store.UpdateProduct(ctx, product)
	if err != nil {
		errResp := errors.ErrResponse{
			Code:      errors.ErrorBadRequest,
			MessageID: "",
			Message:   "Failed to update product",
			Reason:    err.Error(),
		}
		return core.Product{}, errors.Errorf(errResp)
	}

	return product, nil
}

func (p *ProductService) UpdateProductQuantity(ctx context.Context, req core.UpdateProductQuantityRequest) (core.Product, error) {
	product, err := p.Store.GetProduct(ctx, req.ID)
	if err != nil {
		errResp := errors.ErrResponse{
			Code:      errors.ErrorNotFound,
			MessageID: "",
			Message:   "Product not found",
			Reason:    err.Error(),
		}
		return core.Product{}, errors.Errorf(errResp)
	}

	product.QuantityInStock += req.NewQuantity

	product, err = p.Store.UpdateProduct(ctx, product)
	if err != nil {
		errResp := errors.ErrResponse{
			Code:      errors.ErrorBadRequest,
			MessageID: "",
			Message:   "Failed to update product",
			Reason:    err.Error(),
		}
		return core.Product{}, errors.Errorf(errResp)
	}

	return product, nil
}

func (p *ProductService) GetProductByID(ctx context.Context, productID uuid.UUID) (core.Product, error) {
	product, err := p.Store.GetProduct(ctx, productID)
	if err != nil {
		errResp := errors.ErrResponse{
			Code:      errors.ErrorNotFound,
			MessageID: "",
			Message:   "Product not found",
			Reason:    err.Error(),
		}
		return core.Product{}, errors.Errorf(errResp)
	}

	return product, nil
}

func (p *ProductService) ListProducts(ctx context.Context, userID uuid.UUID) ([]core.Product, error) {
	products, err := p.Store.ListProducts(ctx, userID)
	if err != nil {
		errResp := errors.ErrResponse{
			Code:      errors.ErrorBadRequest,
			MessageID: "",
			Message:   "Failed to list products",
			Reason:    err.Error(),
		}
		return nil, errors.Errorf(errResp)
	}

	return products, nil
}

func (p *ProductService) DeleteProduct(ctx context.Context, req core.DeleteProductRequest) error {
	product, err := p.Store.DeleteProduct(ctx, req)
	if err != nil {
		errResp := errors.ErrResponse{
			Code:      errors.ErrorBadRequest,
			MessageID: "",
			Message:   "Failed to delete product",
			Reason:    err.Error(),
		}
		return errors.Errorf(errResp)
	}

	go func() {
		_ = p.s3Store.DeleteFile(product.SKU)
	}()

	return nil
}

func (p *ProductService) GetProductSettings(ctx context.Context, userID uuid.UUID) (core.ProductSettings, error) {
	settings, err := p.Store.GetProductSettings(ctx, userID)
	if err != nil {
		errResp := errors.ErrResponse{
			Code:      errors.ErrorNotFound,
			MessageID: "",
			Message:   "Product settings not found",
			Reason:    err.Error(),
		}
		return core.ProductSettings{}, errors.Errorf(errResp)
	}

	return settings, nil
}

func (p *ProductService) UpdateProductSettings(ctx context.Context, req core.ProductSettings) (core.ProductSettings, error) {
	settings, err := p.Store.UpdateProductSettings(ctx, req)
	if err != nil {
		errResp := errors.ErrResponse{
			Code:      errors.ErrorBadRequest,
			MessageID: "",
			Message:   "Failed to update product settings",
			Reason:    err.Error(),
		}
		return core.ProductSettings{}, errors.Errorf(errResp)
	}

	return settings, nil
}
