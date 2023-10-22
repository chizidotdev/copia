package service

import (
	"github.com/chizidotdev/copia/repository"
)

type Service struct {
	OrderService
}

func NewService(store *repository.Repository) *Service {
	order := NewOrderService(store)

	return &Service{
		OrderService: order,
	}
}
