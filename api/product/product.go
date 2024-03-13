package product

import (
	"github.com/chizidotdev/shop/repository"
	"github.com/chizidotdev/shop/repository/adapters"
)

const (
	storeIDParam   = "storeID"
	productIDParam = "productID"
)

type ProductHandler struct {
	pgStore *repository.Repository
	s3Store *adapters.S3Store
}

func NewProductHandler(pgStore *repository.Repository) *ProductHandler {
	return &ProductHandler{pgStore: pgStore}
}
