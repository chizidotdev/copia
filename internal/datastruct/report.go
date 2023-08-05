package datastruct

import (
	"time"
)

type ReportRow struct {
	TotalItems    int64 `json:"total_items"`
	LowStockItems int64 `json:"low_stock_items"`
	RecentSales   int64 `json:"recent_sales"`
	//PendingOrders    int64                `json:"pending_orders"`
	SalesPerformance float64              `json:"sales_performance"`
	PriceSoldByDate  []PriceSoldByDateRow `json:"price_sold_by_date"`
	PriceSoldByWeek  []PriceSoldByWeekRow `json:"price_sold_by_week"`
}

type PriceSoldByDateRow struct {
	Date           time.Time `json:"date"`
	TotalSalePrice int64     `json:"total_sale_price"`
}

type PriceSoldByWeekRow struct {
	Date           time.Time `json:"date"`
	TotalSalePrice int64     `json:"total_sale_price"`
}
