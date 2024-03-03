package store

import (
	"github.com/chizidotdev/shop/repository"
)

type StoreHandler struct {
	pgStore *repository.Repository
}

func NewStoreHandler(pgStore *repository.Repository) *StoreHandler {
	return &StoreHandler{pgStore: pgStore}
}
