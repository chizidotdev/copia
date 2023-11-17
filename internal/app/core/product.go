package core

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"strings"
	"time"

	"github.com/google/uuid"
)

type ProductRepository interface {
	ListProducts(ctx context.Context, userID uuid.UUID) ([]Product, error)
	GetProduct(ctx context.Context, id uuid.UUID) (Product, error)
	CreateProduct(ctx context.Context, arg Product) (Product, error)
	UpdateProduct(ctx context.Context, arg Product) (Product, error)
	DeleteProduct(ctx context.Context, arg DeleteProductRequest) (Product, error)
	UpdateProductSettings(ctx context.Context, arg ProductSettings) (ProductSettings, error)
	GetProductSettings(ctx context.Context, userID uuid.UUID) (ProductSettings, error)
}

type Product struct {
	ID              uuid.UUID `json:"ID"`
	UserID          uuid.UUID `json:"userID"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	Price           float32   `json:"price"`
	QuantityInStock int       `json:"quantityInStock"`
	ImageURL        string    `json:"imageURL"`
	SKU             string    `json:"SKU"`
}

type ProductSettings struct {
	UserID       uuid.UUID `json:"userID"`
	ReorderPoint int       `json:"reorderPoint" binding:"required,min=0"`
}

type ProductRequest struct {
	UserID          uuid.UUID
	Name            string                `form:"name" binding:"required"`
	Description     string                `form:"description"`
	Price           float32               `form:"price" binding:"required"`
	QuantityInStock int                   `form:"quantityInStock" binding:"required"`
	Image           *multipart.FileHeader `form:"image" binding:"required"`
}

type ProductImageRequest struct {
	ID     uuid.UUID
	UserID uuid.UUID
	Image  *multipart.FileHeader `form:"image" binding:"required"`
}

type UpdateProductQuantityRequest struct {
	ID          uuid.UUID `json:"ID"`
	UserID      uuid.UUID `json:"userID"`
	NewQuantity int       `json:"newQuantity" binding:"required"`
}

type DeleteProductRequest struct {
	UserID uuid.UUID `json:"userID"`
	ID     uuid.UUID `json:"ID"`
}

func GenerateSKU(ID, Name string) string {
	randomGen := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomString := make([]byte, 4)
	chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	for i := range randomString {
		randomString[i] = chars[randomGen.Intn(len(chars))]
	}

	left := strings.ToUpper(Name[:2])
	middle := strings.ToUpper(ID[:4])
	right := string(randomString)

	return fmt.Sprintf("%s-%s-%s", left, middle, right)
}

func ParseImage(image *multipart.FileHeader) (io.Reader, error) {
	//file := imgString[strings.IndexByte(imgString, ',')+1:]

	return image.Open()
}
