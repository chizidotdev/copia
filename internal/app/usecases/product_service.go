package usecases

import (
	"context"
	"github.com/chizidotdev/copia/internal/app/core"
	"github.com/chizidotdev/copia/pkg/errors"
	"github.com/google/uuid"
)

type ProductRepository interface {
	ListProducts(ctx context.Context, userID uuid.UUID) ([]core.Product, error)
	GetProduct(ctx context.Context, id uuid.UUID) (core.Product, error)
	CreateProduct(ctx context.Context, arg core.Product) (core.Product, error)
	UpdateProduct(ctx context.Context, arg core.Product) (core.Product, error)
	DeleteProduct(ctx context.Context, arg core.DeleteProductRequest) error
}

type ProductService struct {
	Store ProductRepository
}

func NewProductService(productRepo ProductRepository) *ProductService {
	return &ProductService{
		Store: productRepo,
	}
}

func (p *ProductService) CreateProduct(ctx context.Context, req core.ProductRequest) (core.Product, error) {
	sku := "VA-SU-123"

	product, err := p.Store.CreateProduct(ctx, core.Product{
		UserID:          req.UserID,
		Name:            req.Name,
		Description:     req.Description,
		Price:           req.Price,
		QuantityInStock: req.QuantityInStock,
		ImageURL:        req.ImageURL,
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
	err := p.Store.DeleteProduct(ctx, req)
	if err != nil {
		return errors.Errorf(errors.ErrorBadRequest, "Failed to delete product: "+err.Error())
	}

	return nil
}
