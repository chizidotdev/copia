package core

import (
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"math/rand"
	"strings"
	"time"
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
	Image           string    `json:"image" binding:"required"`
}

type ProductImageRequest struct {
	ID     uuid.UUID `json:"ID"`
	UserID uuid.UUID `json:"userID"`
	Image  string    `json:"image" binding:"required"`
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

func ParseImage(imgString string) ([]byte, error) {
	file := imgString[strings.IndexByte(imgString, ',')+1:]

	return base64.StdEncoding.DecodeString(file)
}
