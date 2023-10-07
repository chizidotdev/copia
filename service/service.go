package service

import (
	"github.com/chizidotdev/copia/repository"
	"github.com/chizidotdev/copia/util"
)

type Service struct {
	OrderService
	UserService
	TokenManager
	*AuthService
}

func NewService(store *repository.Repository) *Service {
	order := NewOrderService(store)
	user := NewUserService(store)
	tokenManager := NewTokenManagerService(util.EnvVars.AuthSecret)
	auth := NewAuthenticator()

	return &Service{
		OrderService: order,
		UserService:  user,
		TokenManager: tokenManager,
		AuthService:  auth,
	}
}
