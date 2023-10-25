package core

import (
	"github.com/google/uuid"
)

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

type ProductRequest struct {
	UserID          uuid.UUID `json:"userID"`
	Name            string    `json:"name" binding:"required"`
	Description     string    `json:"description"`
	Price           float32   `json:"price" binding:"required"`
	QuantityInStock int       `json:"quantityInStock" binding:"required"`
	ImageURL        string    `json:"imageURL"`
}

type DeleteProductRequest struct {
	UserID uuid.UUID `json:"userID"`
	ID     uuid.UUID `json:"ID"`
}
