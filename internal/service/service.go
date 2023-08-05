package service

import (
	"github.com/chizidotdev/copia/internal/repository"
)

type Service struct {
	OrderService
}

func NewService(store *repository.Store) *Service {
	order := NewOrderService(store)

	return &Service{
		OrderService: order,
	}
}
