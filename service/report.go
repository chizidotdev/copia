package service

import (
	"context"
	"fmt"
	"github.com/chizidotdev/copia/repository"
	"github.com/chizidotdev/copia/util"
	"math"
)

type ReportService interface {
	GetReport(ctx context.Context, userEmail string) (ReportRow, error)
}

type reportService struct {
	Store *repository.Repository
}

func NewDashboardService(store *repository.Repository) ReportService {
	return &reportService{
		Store: store,
	}
}

type ReportRow struct {
	TotalItems    int64 `json:"total_items"`
	LowStockItems int64 `json:"low_stock_items"`
	RecentSales   int64 `json:"recent_sales"`
	//PendingOrders    int64                `json:"pending_orders"`
	SalesPerformance float64                         `json:"sales_performance"`
	PriceSoldByDate  []repository.PriceSoldByDateRow `json:"price_sold_by_date"`
	PriceSoldByWeek  []repository.PriceSoldByWeekRow `json:"price_sold_by_week"`
}

func (r *reportService) GetReport(ctx context.Context, userEmail string) (ReportRow, error) {
	inventory, err := r.Store.GetInventoryStats(ctx, userEmail)
	if err != nil {
		errMsg := fmt.Errorf("failed to get inventory stats: %w", err)
		return ReportRow{}, errMsg
	}

	salesPerformance, err := r.getSalesPerformance(ctx, userEmail)
	if err != nil {
		return ReportRow{}, err
	}

	priceSoldByDate, err := r.getPriceSoldByDate(ctx, userEmail)
	if err != nil {
		return ReportRow{}, err
	}

	priceSoldByWeek, err := r.getPriceSoldByWeek(ctx, userEmail)
	if err != nil {
		return ReportRow{}, err
	}

	report := ReportRow{
		TotalItems:    inventory.TotalItems,
		LowStockItems: inventory.LowStockItems,
		RecentSales:   inventory.RecentSales,
		//PendingOrders:    inventory.PendingOrders,
		SalesPerformance: salesPerformance,
		PriceSoldByDate:  priceSoldByDate,
		PriceSoldByWeek:  priceSoldByWeek,
	}

	return report, nil
}

func (r *reportService) getSalesPerformance(ctx context.Context, userEmail string) (float64, error) {
	currWeekSale, err := r.Store.CurrentWeekSales(ctx, userEmail)
	if err != nil {
		errMsg := fmt.Errorf("failed to get current week sales: %w", err)
		return 0, errMsg
	}

	lastWeekSales, err := r.Store.LastWeekSales(ctx, userEmail)
	if err != nil {
		errMsg := fmt.Errorf("failed to get last week sales: %w", err)
		return 0, errMsg
	}

	var salesPerformance float64
	if lastWeekSales == 0 {
		if currWeekSale == 0 {
			salesPerformance = 0
		} else {
			salesPerformance = 100
		}
	} else {
		diff := util.CalcPercentageDiff(lastWeekSales, currWeekSale)
		salesPerformance = math.Floor(diff)
	}

	return salesPerformance, nil
}

func (r *reportService) getPriceSoldByDate(ctx context.Context, userEmail string) ([]repository.PriceSoldByDateRow, error) {
	priceSoldByDate, err := r.Store.PriceSoldByDate(ctx, userEmail)
	if err != nil {
		errMsg := fmt.Errorf("failed to get price sold by date: %w", err)
		return nil, errMsg
	}

	return priceSoldByDate, nil
}

func (r *reportService) getPriceSoldByWeek(ctx context.Context, userEmail string) ([]repository.PriceSoldByWeekRow, error) {
	priceSoldByWeek, err := r.Store.PriceSoldByWeek(ctx, userEmail)
	if err != nil {
		errMsg := fmt.Errorf("failed to get price sold by week: %w", err)
		return nil, errMsg
	}

	return priceSoldByWeek, nil
}
