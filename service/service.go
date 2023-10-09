package service

import (
	"github.com/chizidotdev/copia/repository"
)

type Service struct {
	OrderService
	UserService
	*AuthService
}

func NewService(store *repository.Repository) *Service {
	order := NewOrderService(store)
	user := NewUserService(store)
	auth := NewAuthenticator()

	return &Service{
		OrderService: order,
		UserService:  user,
		AuthService:  auth,
	}
}
