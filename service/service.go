package service

import (
	"github.com/chizidotdev/copia/repository"
	"github.com/chizidotdev/copia/util"
)

type Service struct {
	OrderService
	TokenManager
}

func NewService(store *repository.Repository) *Service {
	order := NewOrderService(store)
	tokenManager := NewTokenManagerService(util.EnvVars.AuthSecret)

	return &Service{
		OrderService: order,
		TokenManager: tokenManager,
	}
}
