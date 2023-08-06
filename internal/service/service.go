package service

import (
	"github.com/chizidotdev/copia/internal/repository"
	"github.com/chizidotdev/copia/pkg/utils"
)

type Service struct {
	OrderService
	TokenManager
}

func NewService(store *repository.Store) *Service {
	order := NewOrderService(store)
	tokenManager := NewTokenManagerService(utils.EnvVars.AuthSecret)

	return &Service{
		OrderService: order,
		TokenManager: tokenManager,
	}
}
