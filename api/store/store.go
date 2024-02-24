package store

import (
	"github.com/chizidotdev/shop/repository"
)

type StoreHandler struct {
	pgStore *repository.Queries
}

func NewStoreHandler(pgStore *repository.Queries) *StoreHandler {
	return &StoreHandler{pgStore: pgStore}
}
