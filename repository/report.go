package repository

import (
	"context"
	"time"

	"gorm.io/gorm/clause"
)

type GetInventoryStatsRow struct {
	TotalItems    int64 `json:"total_items"`
	LowStockItems int64 `json:"low_stock_items"`
	RecentSales   int64 `json:"recent_sales"`
	//PendingOrders int64 `json:"pending_orders"`
}

func (r *Repository) GetInventoryStats(_ context.Context, userEmail string) (GetInventoryStatsRow, error) {
	var result GetInventoryStatsRow
	err := r.DB.Model(&Order{}).Where("user_email = ?", userEmail).Count(&result.TotalItems).Error
	if err != nil {
		return result, err
	}

	err = r.DB.Model(&Order{}).Where("user_email = ? AND quantity < 5", userEmail).Count(&result.LowStockItems).Error
	if err != nil {
		return result, err
	}

	err = r.DB.Model(&Customer{}).Where("user_email = ?", userEmail).Count(&result.RecentSales).Error
	if err != nil {
		return result, err
	}

	return result, err
}

type PriceSoldByWeekRow struct {
	Date           time.Time `json:"date"`
	TotalSalePrice int64     `json:"total_sale_price"`
}

func (r *Repository) PriceSoldByWeek(_ context.Context, userEmail string) ([]PriceSoldByWeekRow, error) {
	var results []PriceSoldByWeekRow
	err := r.DB.Model(&Customer{}).
		Select("DATE_TRUNC('week', sale_date) AS date, SUM(sale_price * quantity_sold) AS total_sale_price").
		Where("user_email = ?", userEmail).
		Group("DATE_TRUNC('week', sale_date)").
		Order(clause.Expr{SQL: "DATE_TRUNC('week', sale_date)"}).
		Scan(&results).Error

	return results, err
}

type PriceSoldByDateRow struct {
	Date           time.Time `json:"date"`
	TotalSalePrice int64     `json:"total_sale_price"`
}

func (r *Repository) PriceSoldByDate(_ context.Context, userEmail string) ([]PriceSoldByDateRow, error) {
	var results []PriceSoldByDateRow
	err := r.DB.Model(&Customer{}).
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

func (r *Repository) CurrentWeekSales(_ context.Context, userEmail string) (int32, error) {
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

func (r *Repository) LastWeekSales(_ context.Context, userEmail string) (int32, error) {
	return 0, nil
}
