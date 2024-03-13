package order

import (
	"github.com/chizidotdev/shop/repository"
)

const (
	orderIDParam     = "orderID"
	orderItemIDParam = "orderItemID"
)

type OrderHandler struct {
	pgStore *repository.Repository
}

func NewOrderHandler(pgStore *repository.Repository) *OrderHandler {
	return &OrderHandler{
		pgStore: pgStore,
	}
}
