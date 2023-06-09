package service

import (
	"context"
	"fmt"
	"github.com/chizidotdev/copia/internal/datastruct"
	"github.com/chizidotdev/copia/internal/repository"
)

type SaleService interface {
	ListSalesByUser(ctx context.Context, user datastruct.UserJWT) ([]repository.ListSalesByUserIdRow, error)
	ListSalesByItem(ctx context.Context, req repository.ListSalesParams) ([]repository.ListSalesRow, error)
	CreateSale(ctx context.Context, req repository.CreateSaleParams) (repository.Sale, error)
	UpdateSale(ctx context.Context, req repository.UpdateSaleParams) (repository.Sale, error)
	GetSaleByID(ctx context.Context, params repository.GetSaleParams) (repository.Sale, error)
	DeleteSale(ctx context.Context, req repository.DeleteSaleParams) error
}

type saleService struct {
	Store *repository.Store
}

func NewSaleService(store *repository.Store) SaleService {
	return &saleService{
		Store: store,
	}
}

func (s *saleService) ListSalesByUser(ctx context.Context, user datastruct.UserJWT) ([]repository.ListSalesByUserIdRow, error) {
	sales, err := s.Store.ListSalesByUserId(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return sales, nil
}

func (s *saleService) ListSalesByItem(ctx context.Context, req repository.ListSalesParams) ([]repository.ListSalesRow, error) {
	_, err := s.Store.GetItem(ctx, req.ItemID)
	if err != nil {
		return nil, err
	}

	sales, err := s.Store.ListSales(ctx, req)
	if err != nil {
		return nil, err
	}

	return sales, nil
}

func (s *saleService) CreateSale(ctx context.Context, req repository.CreateSaleParams) (repository.Sale, error) {
	var sale repository.Sale
	err := s.Store.ExecTx(ctx, func(query *repository.Queries) error {
		var err error
		sale, err = query.CreateSale(ctx, req)
		if err != nil {
			return err
		}
		_, err = query.UpdateItemQuantity(ctx, repository.UpdateItemQuantityParams{
			ID:       req.ItemID,
			Quantity: -req.QuantitySold,
			UserID:   req.UserID,
		})
		return err
	})
	if err != nil {
		return repository.Sale{}, err
	}

	return sale, nil
}

func (s *saleService) UpdateSale(ctx context.Context, req repository.UpdateSaleParams) (repository.Sale, error) {
	initialSale, err := s.Store.GetSale(ctx, repository.GetSaleParams{
		ID:     req.ID,
		UserID: req.UserID,
	})
	if err != nil {
		errMessage := fmt.Errorf("sale not found: %v", err)
		return repository.Sale{}, errMessage
	}

	quantityDiff := req.QuantitySold - initialSale.QuantitySold

	args := repository.UpdateSaleParams{
		ID:           req.ID,
		QuantitySold: req.QuantitySold,
		SalePrice:    req.SalePrice,
		CustomerName: req.CustomerName,
		SaleDate:     req.SaleDate,
		UserID:       req.UserID,
	}

	itemArgs := repository.UpdateItemQuantityParams{
		ID:       initialSale.ItemID,
		Quantity: -quantityDiff,
		UserID:   req.UserID,
	}

	var sale repository.Sale
	err = s.Store.ExecTx(ctx, func(query *repository.Queries) error {
		var err error
		sale, err = query.UpdateSale(ctx, args)
		if err != nil {
			return fmt.Errorf("error updating sale: %w", err)
		}
		_, err = query.UpdateItemQuantity(ctx, itemArgs)
		if err != nil {
			return fmt.Errorf("error updating item quantity: %w", err)
		}
		return err
	})
	if err != nil {
		return repository.Sale{}, err
	}

	return sale, nil
}

func (s *saleService) GetSaleByID(ctx context.Context, req repository.GetSaleParams) (repository.Sale, error) {
	sale, err := s.Store.GetSale(ctx, req)
	if err != nil {
		return repository.Sale{}, err
	}

	return sale, nil
}

func (s *saleService) DeleteSale(ctx context.Context, req repository.DeleteSaleParams) error {
	sale, err := s.Store.GetSale(ctx, repository.GetSaleParams{
		ID:     req.ID,
		UserID: req.UserID,
	})
	if err != nil {
		return fmt.Errorf("sale not found: %v", err)
	}

	err = s.Store.ExecTx(ctx, func(query *repository.Queries) error {
		err := query.DeleteSale(ctx, repository.DeleteSaleParams{
			ID:     req.ID,
			UserID: req.UserID,
		})
		if err != nil {
			return fmt.Errorf("an error occurred while deleting sale: %v", err)
		}

		_, err = query.UpdateItemQuantity(ctx, repository.UpdateItemQuantityParams{
			ID:       sale.ItemID,
			Quantity: sale.QuantitySold,
			UserID:   sale.UserID,
		})
		if err != nil {
			return fmt.Errorf("an error occurred while updating sale item quantity: %v", err)
		}

		return err
	})
	if err != nil {
		return err
	}

	return nil
}
