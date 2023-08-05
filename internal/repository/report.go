package repository

import (
	"context"

	"github.com/chizidotdev/copia/internal/datastruct"
	"gorm.io/gorm/clause"
)

type GetInventoryStatsRow struct {
	TotalItems    int64 `json:"total_items"`
	LowStockItems int64 `json:"low_stock_items"`
	RecentSales   int64 `json:"recent_sales"`
	//PendingOrders int64 `json:"pending_orders"`
}

func (s *Store) GetInventoryStats(_ context.Context, userEmail string) (GetInventoryStatsRow, error) {
	var result GetInventoryStatsRow
	err := s.DB.Model(&Order{}).Where("user_email = ?", userEmail).Count(&result.TotalItems).Error
	if err != nil {
		return result, err
	}

	err = s.DB.Model(&Order{}).Where("user_email = ? AND quantity < 5", userEmail).Count(&result.LowStockItems).Error
	if err != nil {
		return result, err
	}

	err = s.DB.Model(&Customer{}).Where("user_email = ?", userEmail).Count(&result.RecentSales).Error
	if err != nil {
		return result, err
	}

	return result, err
}

func (s *Store) PriceSoldByWeek(_ context.Context, userEmail string) ([]datastruct.PriceSoldByWeekRow, error) {
	var results []datastruct.PriceSoldByWeekRow
	err := s.DB.Model(&Customer{}).
		Select("DATE_TRUNC('week', sale_date) AS date, SUM(sale_price * quantity_sold) AS total_sale_price").
		Where("user_email = ?", userEmail).
		Group("DATE_TRUNC('week', sale_date)").
		Order(clause.Expr{SQL: "DATE_TRUNC('week', sale_date)"}).
		Scan(&results).Error

	return results, err
}

func (s *Store) PriceSoldByDate(_ context.Context, userEmail string) ([]datastruct.PriceSoldByDateRow, error) {
	var results []datastruct.PriceSoldByDateRow
	err := s.DB.Model(&Customer{}).
		Select("DATE_TRUNC('day', sale_date) AS date, SUM(sale_price * quantity_sold) AS total_sale_price").
		Where("user_email = ?", userEmail).
		Group("DATE_TRUNC('day', sale_date)").
		Order(clause.Expr{SQL: "DATE_TRUNC('day', sale_date)"}).
		Scan(&results).Error

	return results, err
}

//const currentWeekSales = `-- name: CurrentWeekSales :one
//SELECT CAST(COALESCE(SUM(quantity_sold), 0) AS INTEGER) AS total_quantity_sold
//FROM sales
//WHERE sale_date >= DATE_TRUNC('week', CURRENT_DATE)
//  AND sale_date < DATE_TRUNC('week', CURRENT_DATE) + INTERVAL '1 week'
//  AND user_id = $1
//`

func (s *Store) CurrentWeekSales(_ context.Context, userEmail string) (int32, error) {
	return 0, nil
}

//const lastWeekSales = `-- name: LastWeekSales :one
//SELECT CAST(COALESCE(SUM(quantity_sold), 0) AS INTEGER) AS total_quantity_sold
//FROM sales
//WHERE sale_date >= DATE_TRUNC('week', CURRENT_DATE) - INTERVAL '1 week'
//  AND sale_date
//    < DATE_TRUNC('week'
//    , CURRENT_DATE)
//  AND user_id = $1
//`

func (s *Store) LastWeekSales(_ context.Context, userEmail string) (int32, error) {
	return 0, nil
}
