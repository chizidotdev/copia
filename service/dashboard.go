package service

import (
	"context"
	"fmt"
	repository2 "github.com/chizidotdev/copia/repository"
	"github.com/chizidotdev/copia/util"
	"math"
)

type DashboardService interface {
	GetDashboard(ctx context.Context, userEmail string) (ReportRow, error)
}

type dashboardService struct {
	Store *repository2.Store
}

func NewDashboardService(store *repository2.Store) DashboardService {
	return &dashboardService{
		Store: store,
	}
}

type ReportRow struct {
	TotalItems    int64 `json:"total_items"`
	LowStockItems int64 `json:"low_stock_items"`
	RecentSales   int64 `json:"recent_sales"`
	//PendingOrders    int64                `json:"pending_orders"`
	SalesPerformance float64                          `json:"sales_performance"`
	PriceSoldByDate  []repository2.PriceSoldByDateRow `json:"price_sold_by_date"`
	PriceSoldByWeek  []repository2.PriceSoldByWeekRow `json:"price_sold_by_week"`
}

func (d *dashboardService) GetDashboard(ctx context.Context, userEmail string) (ReportRow, error) {
	inventory, err := d.Store.GetInventoryStats(ctx, userEmail)
	if err != nil {
		errMsg := fmt.Errorf("failed to get inventory stats: %w", err)
		return ReportRow{}, errMsg
	}

	salesPerformance, err := d.getSalesPerformance(ctx, userEmail)
	if err != nil {
		return ReportRow{}, err
	}

	priceSoldByDate, err := d.getPriceSoldByDate(ctx, userEmail)
	if err != nil {
		return ReportRow{}, err
	}

	priceSoldByWeek, err := d.getPriceSoldByWeek(ctx, userEmail)
	if err != nil {
		return ReportRow{}, err
	}

	dashboard := ReportRow{
		TotalItems:    inventory.TotalItems,
		LowStockItems: inventory.LowStockItems,
		RecentSales:   inventory.RecentSales,
		//PendingOrders:    inventory.PendingOrders,
		SalesPerformance: salesPerformance,
		PriceSoldByDate:  priceSoldByDate,
		PriceSoldByWeek:  priceSoldByWeek,
	}

	return dashboard, nil
}

func (d *dashboardService) getSalesPerformance(ctx context.Context, userEmail string) (float64, error) {
	currWeekSale, err := d.Store.CurrentWeekSales(ctx, userEmail)
	if err != nil {
		errMsg := fmt.Errorf("failed to get current week sales: %w", err)
		return 0, errMsg
	}

	lastWeekSales, err := d.Store.LastWeekSales(ctx, userEmail)
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

func (d *dashboardService) getPriceSoldByDate(ctx context.Context, userEmail string) ([]repository2.PriceSoldByDateRow, error) {
	priceSoldByDate, err := d.Store.PriceSoldByDate(ctx, userEmail)
	if err != nil {
		errMsg := fmt.Errorf("failed to get price sold by date: %w", err)
		return nil, errMsg
	}

	return priceSoldByDate, nil
}

func (d *dashboardService) getPriceSoldByWeek(ctx context.Context, userEmail string) ([]repository2.PriceSoldByWeekRow, error) {
	priceSoldByWeek, err := d.Store.PriceSoldByWeek(ctx, userEmail)
	if err != nil {
		errMsg := fmt.Errorf("failed to get price sold by week: %w", err)
		return nil, errMsg
	}

	return priceSoldByWeek, nil
}
