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

	image, err := core.ParseImage(req.Image)
	if err != nil {
		return core.Product{}, errors.Errorf(errors.ErrorBadRequest, "Failed to open file: "+err.Error())
	}

	imageUrl, err := p.s3Store.UploadFile(sku, image)
	if err != nil {
		return core.Product{}, errors.Errorf(errors.ErrorBadRequest, "Failed to upload file: "+err.Error())
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
		return core.Product{}, errors.Errorf(errors.ErrorBadRequest, "Failed to create product: "+err.Error())
	}

	return product, nil
}

func (p *ProductService) UpdateProduct(ctx context.Context, req core.Product) (core.Product, error) {
	product, err := p.Store.UpdateProduct(ctx, req)
	if err != nil {
		return core.Product{}, errors.Errorf(errors.ErrorBadRequest, "Failed to update product: "+err.Error())
	}

	return product, nil
}

func (p *ProductService) UpdateProductImage(ctx context.Context, req core.ProductImageRequest) (core.Product, error) {
	product, err := p.Store.GetProduct(ctx, req.ID)
	if err != nil {
		return core.Product{}, errors.Errorf(errors.ErrorNotFound, "Product not found")
	}

	image, err := core.ParseImage(req.Image)
	if err != nil {
		return core.Product{}, errors.Errorf(errors.ErrorBadRequest, "Failed to open file: "+err.Error())
	}

	imageUrl, err := p.s3Store.UploadFile(product.SKU, image)
	if err != nil {
		return core.Product{}, errors.Errorf(errors.ErrorBadRequest, "Failed to upload file: "+err.Error())
	}

	product.ImageURL = imageUrl

	product, err = p.Store.UpdateProduct(ctx, product)
	if err != nil {
		return core.Product{}, errors.Errorf(errors.ErrorBadRequest, "Failed to update product: "+err.Error())
	}

	return product, nil
}

func (p *ProductService) UpdateProductQuantity(ctx context.Context, req core.UpdateProductQuantityRequest) (core.Product, error) {
	product, err := p.Store.GetProduct(ctx, req.ID)
	if err != nil {
		return core.Product{}, errors.Errorf(errors.ErrorNotFound, "Product not found")
	}

	product.QuantityInStock += req.NewQuantity

	product, err = p.Store.UpdateProduct(ctx, product)
	if err != nil {
		return core.Product{}, errors.Errorf(errors.ErrorBadRequest, "Failed to update product: "+err.Error())
	}

	return product, nil
}

func (p *ProductService) GetProductByID(ctx context.Context, productID uuid.UUID) (core.Product, error) {
	product, err := p.Store.GetProduct(ctx, productID)
	if err != nil {
		return core.Product{}, errors.Errorf(errors.ErrorNotFound, "Product not found: "+err.Error())
	}

	return product, nil
}

func (p *ProductService) ListProducts(ctx context.Context, userID uuid.UUID) ([]core.Product, error) {
	products, err := p.Store.ListProducts(ctx, userID)
	if err != nil {
		return nil, errors.Errorf(errors.ErrorNotFound, "Products not found: "+err.Error())
	}

	return products, nil
}

func (p *ProductService) DeleteProduct(ctx context.Context, req core.DeleteProductRequest) error {
	product, err := p.Store.DeleteProduct(ctx, req)
	if err != nil {
		return errors.Errorf(errors.ErrorBadRequest, "Failed to delete product: "+err.Error())
	}

	go func() {
		_ = p.s3Store.DeleteFile(product.SKU)
	}()

	return nil
}

func (p *ProductService) UpdateProductSettings(ctx context.Context, req core.ProductSettings) (core.ProductSettings, error) {
	settings, err := p.Store.UpdateProductSettings(ctx, req)
	if err != nil {
		return core.ProductSettings{}, errors.Errorf(errors.ErrorBadRequest, "Failed to update product settings: "+err.Error())
	}

	return settings, nil
}

func (p *ProductService) GetProductSettings(ctx context.Context, userID uuid.UUID) (core.ProductSettings, error) {
	settings, err := p.Store.GetProductSettings(ctx, userID)
	if err != nil {
		return core.ProductSettings{}, errors.Errorf(errors.ErrorNotFound, "Product settings not found: "+err.Error())
	}

	return settings, nil
}
