package adapters

import (
	"context"
	"github.com/chizidotdev/copia/internal/app/core"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
)

type Product struct {
	Base
	UserID          uuid.UUID `gorm:"not null;uniqueIndex" json:"userID"`
	Name            string    `gorm:"not null" json:"name"`
	Description     string    `gorm:"not null" json:"description"`
	Price           float32   `gorm:"not null" json:"price"`
	QuantityInStock int       `gorm:"not null" json:"quantityInStock"`
	ImageURL        string    `gorm:"not null" json:"imageURL"`
	SKU             string    `gorm:"not null" json:"SKU"`
}

type ProductSettings struct {
	UserID       uuid.UUID `gorm:"not null;uniqueIndex" json:"userID"`
	ReorderPoint int       `gorm:"not null;default:0" json:"reorderPoint"`
}

var _ core.ProductRepository = (*ProductRepositoryImpl)(nil)

type ProductRepositoryImpl struct {
	DB *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepositoryImpl {
	err := db.AutoMigrate(&Product{}, &ProductSettings{})
	if err != nil {
		log.Panic("Failed to migrate Product", err)
	}
	return &ProductRepositoryImpl{DB: db}
}

func (p *ProductRepositoryImpl) ListProducts(_ context.Context, userID uuid.UUID) ([]core.Product, error) {
	var products []core.Product
	result := p.DB.Find(&products, "user_id = ?", userID)
	return products, result.Error
}

func (p *ProductRepositoryImpl) GetProduct(_ context.Context, id uuid.UUID) (core.Product, error) {
	var product core.Product
	result := p.DB.First(&product, "id = ?", id)
	return product, result.Error
}

func (p *ProductRepositoryImpl) CreateProduct(_ context.Context, arg core.Product) (core.Product, error) {
	product := Product{
		UserID:          arg.UserID,
		Name:            arg.Name,
		Description:     arg.Description,
		Price:           arg.Price,
		QuantityInStock: arg.QuantityInStock,
		ImageURL:        arg.ImageURL,
		SKU:             arg.SKU,
	}
	result := p.DB.Create(&product)
	return core.Product{
		ID:              product.ID,
		UserID:          product.UserID,
		Name:            product.Name,
		Description:     product.Description,
		Price:           product.Price,
		QuantityInStock: product.QuantityInStock,
		ImageURL:        product.ImageURL,
		SKU:             product.SKU,
	}, result.Error
}

func (p *ProductRepositoryImpl) UpdateProduct(_ context.Context, arg core.Product) (core.Product, error) {
	var product Product
	result := p.DB.Model(&product).
		Clauses(clause.Returning{}).
		Where("id = ? AND user_id = ?", arg.ID, arg.UserID).
		First(&product)
	if result.Error != nil {
		return core.Product{}, result.Error
	}
	err := result.Updates(Product{
		Name:            arg.Name,
		Description:     arg.Description,
		Price:           arg.Price,
		QuantityInStock: arg.QuantityInStock,
		ImageURL:        arg.ImageURL,
		SKU:             arg.SKU,
	}).Error
	if err != nil {
		return core.Product{}, result.Error
	}

	return core.Product{
		ID:              product.ID,
		UserID:          product.UserID,
		Name:            product.Name,
		Description:     product.Description,
		Price:           product.Price,
		QuantityInStock: product.QuantityInStock,
		ImageURL:        product.ImageURL,
		SKU:             product.SKU,
	}, nil
}

func (p *ProductRepositoryImpl) DeleteProduct(_ context.Context, arg core.DeleteProductRequest) (core.Product, error) {
	var product Product
	result := p.DB.First(&product, "id = ? AND user_id = ?", arg.ID, arg.UserID)
	if result.Error != nil {
		return core.Product{}, result.Error
	}
	result.Delete(&product)

	return core.Product{
		ID:              product.ID,
		UserID:          product.UserID,
		Name:            product.Name,
		Description:     product.Description,
		Price:           product.Price,
		QuantityInStock: product.QuantityInStock,
		ImageURL:        product.ImageURL,
		SKU:             product.SKU,
	}, nil
}

func (p *ProductRepositoryImpl) UpdateProductSettings(_ context.Context, arg core.ProductSettings) (core.ProductSettings, error) {
	productSettings := ProductSettings{
		UserID:       arg.UserID,
		ReorderPoint: arg.ReorderPoint,
	}
	result := p.DB.
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}},
			UpdateAll: true,
		}).
		Where("user_id = ?", arg.UserID).
		Create(&productSettings)

	return core.ProductSettings{
		UserID:       productSettings.UserID,
		ReorderPoint: productSettings.ReorderPoint,
	}, result.Error
}
