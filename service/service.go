package service

import (
	"github.com/chizidotdev/copia/repository"
)

type Service struct {
	OrderService
	UserService
}

func NewService(store *repository.Repository) *Service {
	order := NewOrderService(store)
	user := NewUserService(store)

	return &Service{
		OrderService: order,
		UserService:  user,
	}
}
