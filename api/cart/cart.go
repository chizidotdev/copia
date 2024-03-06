package cart

import (
	"github.com/chizidotdev/shop/repository"
)

const (
	cartIDParam = "cartID"
)

type CartHandler struct {
	pgStore *repository.Repository
}

func NewCartHandler(pgStore *repository.Repository) *CartHandler {
	return &CartHandler{
		pgStore: pgStore,
	}
}
