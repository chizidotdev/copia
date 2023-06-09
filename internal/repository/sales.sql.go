// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: sales.sql

package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createSale = `-- name: CreateSale :one
INSERT INTO sales (item_id, user_id, quantity_sold, sale_price, customer_name, sale_date)
VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, item_id, user_id, quantity_sold, sale_price, sale_date, customer_name, created_at, updated_at
`

type CreateSaleParams struct {
	ItemID       uuid.UUID `json:"item_id"`
	UserID       uuid.UUID `json:"user_id"`
	QuantitySold int64     `json:"quantity_sold"`
	SalePrice    float32   `json:"sale_price"`
	CustomerName string    `json:"customer_name"`
	SaleDate     time.Time `json:"sale_date"`
}

func (q *Queries) CreateSale(ctx context.Context, arg CreateSaleParams) (Sale, error) {
	row := q.db.QueryRowContext(ctx, createSale,
		arg.ItemID,
		arg.UserID,
		arg.QuantitySold,
		arg.SalePrice,
		arg.CustomerName,
		arg.SaleDate,
	)
	var i Sale
	err := row.Scan(
		&i.ID,
		&i.ItemID,
		&i.UserID,
		&i.QuantitySold,
		&i.SalePrice,
		&i.SaleDate,
		&i.CustomerName,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const currentWeekSales = `-- name: CurrentWeekSales :one
SELECT CAST(COALESCE(SUM(quantity_sold), 0) AS INTEGER) AS total_quantity_sold
FROM sales
WHERE sale_date >= DATE_TRUNC('week', CURRENT_DATE)
  AND sale_date < DATE_TRUNC('week', CURRENT_DATE) + INTERVAL '1 week'
  AND user_id = $1
`

func (q *Queries) CurrentWeekSales(ctx context.Context, userID uuid.UUID) (int32, error) {
	row := q.db.QueryRowContext(ctx, currentWeekSales, userID)
	var total_quantity_sold int32
	err := row.Scan(&total_quantity_sold)
	return total_quantity_sold, err
}

const deleteSale = `-- name: DeleteSale :exec
DELETE
FROM sales
WHERE (id = $1 AND user_id = $2)
`

type DeleteSaleParams struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
}

func (q *Queries) DeleteSale(ctx context.Context, arg DeleteSaleParams) error {
	_, err := q.db.ExecContext(ctx, deleteSale, arg.ID, arg.UserID)
	return err
}

const getSale = `-- name: GetSale :one
SELECT id, item_id, user_id, quantity_sold, sale_price, sale_date, customer_name, created_at, updated_at
FROM sales
WHERE (id = $1 AND user_id = $2) LIMIT 1
`

type GetSaleParams struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
}

func (q *Queries) GetSale(ctx context.Context, arg GetSaleParams) (Sale, error) {
	row := q.db.QueryRowContext(ctx, getSale, arg.ID, arg.UserID)
	var i Sale
	err := row.Scan(
		&i.ID,
		&i.ItemID,
		&i.UserID,
		&i.QuantitySold,
		&i.SalePrice,
		&i.SaleDate,
		&i.CustomerName,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getSaleForUpdate = `-- name: GetSaleForUpdate :one
SELECT id, item_id, user_id, quantity_sold, sale_price, sale_date, customer_name, created_at, updated_at
FROM sales
WHERE (id = $1 AND user_id = $2) LIMIT 1 FOR NO KEY
UPDATE
`

type GetSaleForUpdateParams struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
}

func (q *Queries) GetSaleForUpdate(ctx context.Context, arg GetSaleForUpdateParams) (Sale, error) {
	row := q.db.QueryRowContext(ctx, getSaleForUpdate, arg.ID, arg.UserID)
	var i Sale
	err := row.Scan(
		&i.ID,
		&i.ItemID,
		&i.UserID,
		&i.QuantitySold,
		&i.SalePrice,
		&i.SaleDate,
		&i.CustomerName,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const lastWeekSales = `-- name: LastWeekSales :one
SELECT CAST(COALESCE(SUM(quantity_sold), 0) AS INTEGER) AS total_quantity_sold
FROM sales
WHERE sale_date >= DATE_TRUNC('week', CURRENT_DATE) - INTERVAL '1 week'
  AND sale_date
    < DATE_TRUNC('week'
    , CURRENT_DATE)
  AND user_id = $1
`

func (q *Queries) LastWeekSales(ctx context.Context, userID uuid.UUID) (int32, error) {
	row := q.db.QueryRowContext(ctx, lastWeekSales, userID)
	var total_quantity_sold int32
	err := row.Scan(&total_quantity_sold)
	return total_quantity_sold, err
}

const listSales = `-- name: ListSales :many
SELECT s.id, s.item_id, s.user_id, s.quantity_sold, s.sale_price, s.sale_date, s.customer_name, s.created_at, s.updated_at, i.title
FROM sales s
         JOIN items i ON s.item_id = i.id
WHERE s.item_id = $1
  AND s.user_id = $2
ORDER BY s.sale_date DESC
`

type ListSalesParams struct {
	ItemID uuid.UUID `json:"item_id"`
	UserID uuid.UUID `json:"user_id"`
}

type ListSalesRow struct {
	ID           uuid.UUID `json:"id"`
	ItemID       uuid.UUID `json:"item_id"`
	UserID       uuid.UUID `json:"user_id"`
	QuantitySold int64     `json:"quantity_sold"`
	SalePrice    float32   `json:"sale_price"`
	SaleDate     time.Time `json:"sale_date"`
	CustomerName string    `json:"customer_name"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Title        string    `json:"title"`
}

func (q *Queries) ListSales(ctx context.Context, arg ListSalesParams) ([]ListSalesRow, error) {
	rows, err := q.db.QueryContext(ctx, listSales, arg.ItemID, arg.UserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListSalesRow{}
	for rows.Next() {
		var i ListSalesRow
		if err := rows.Scan(
			&i.ID,
			&i.ItemID,
			&i.UserID,
			&i.QuantitySold,
			&i.SalePrice,
			&i.SaleDate,
			&i.CustomerName,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Title,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listSalesByUserId = `-- name: ListSalesByUserId :many
SELECT s.id, s.item_id, s.user_id, s.quantity_sold, s.sale_price, s.sale_date, s.customer_name, s.created_at, s.updated_at, i.title
FROM sales s
         JOIN items i ON s.item_id = i.id
WHERE s.user_id = $1
ORDER BY s.sale_date DESC
`

type ListSalesByUserIdRow struct {
	ID           uuid.UUID `json:"id"`
	ItemID       uuid.UUID `json:"item_id"`
	UserID       uuid.UUID `json:"user_id"`
	QuantitySold int64     `json:"quantity_sold"`
	SalePrice    float32   `json:"sale_price"`
	SaleDate     time.Time `json:"sale_date"`
	CustomerName string    `json:"customer_name"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Title        string    `json:"title"`
}

func (q *Queries) ListSalesByUserId(ctx context.Context, userID uuid.UUID) ([]ListSalesByUserIdRow, error) {
	rows, err := q.db.QueryContext(ctx, listSalesByUserId, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListSalesByUserIdRow{}
	for rows.Next() {
		var i ListSalesByUserIdRow
		if err := rows.Scan(
			&i.ID,
			&i.ItemID,
			&i.UserID,
			&i.QuantitySold,
			&i.SalePrice,
			&i.SaleDate,
			&i.CustomerName,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Title,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const priceSoldByDate = `-- name: PriceSoldByDate :many
SELECT DATE_TRUNC('day', sale_date) AS date,
    SUM(sale_price) AS total_sale_price
FROM
    sales
WHERE
    user_id = $1
GROUP BY
    DATE_TRUNC('day', sale_date)
ORDER BY
    DATE_TRUNC('day', sale_date)
`

type PriceSoldByDateRow struct {
	Date           time.Time `json:"date"`
	TotalSalePrice int64 `json:"total_sale_price"`
}

func (q *Queries) PriceSoldByDate(ctx context.Context, userID uuid.UUID) ([]PriceSoldByDateRow, error) {
	rows, err := q.db.QueryContext(ctx, priceSoldByDate, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []PriceSoldByDateRow{}
	for rows.Next() {
		var i PriceSoldByDateRow
		if err := rows.Scan(&i.Date, &i.TotalSalePrice); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const priceSoldByWeek = `-- name: PriceSoldByWeek :many
SELECT DATE_TRUNC('week', sale_date) AS date,
    SUM(sale_price) AS total_sale_price
FROM
    sales
WHERE
    user_id = $1
GROUP BY
    DATE_TRUNC('week', sale_date)
ORDER BY
    DATE_TRUNC('week', sale_date)
`

type PriceSoldByWeekRow struct {
	Date           time.Time `json:"date"`
	TotalSalePrice int64 `json:"total_sale_price"`
}

func (q *Queries) PriceSoldByWeek(ctx context.Context, userID uuid.UUID) ([]PriceSoldByWeekRow, error) {
	rows, err := q.db.QueryContext(ctx, priceSoldByWeek, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []PriceSoldByWeekRow{}
	for rows.Next() {
		var i PriceSoldByWeekRow
		if err := rows.Scan(&i.Date, &i.TotalSalePrice); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateSale = `-- name: UpdateSale :one
UPDATE sales
SET quantity_sold = $2,
    sale_price    = $3,
    customer_name = $4,
    sale_date     = $5
WHERE (id = $1 AND user_id = $6) RETURNING id, item_id, user_id, quantity_sold, sale_price, sale_date, customer_name, created_at, updated_at
`

type UpdateSaleParams struct {
	ID           uuid.UUID `json:"id"`
	QuantitySold int64     `json:"quantity_sold"`
	SalePrice    float32   `json:"sale_price"`
	CustomerName string    `json:"customer_name"`
	SaleDate     time.Time `json:"sale_date"`
	UserID       uuid.UUID `json:"user_id"`
}

func (q *Queries) UpdateSale(ctx context.Context, arg UpdateSaleParams) (Sale, error) {
	row := q.db.QueryRowContext(ctx, updateSale,
		arg.ID,
		arg.QuantitySold,
		arg.SalePrice,
		arg.CustomerName,
		arg.SaleDate,
		arg.UserID,
	)
	var i Sale
	err := row.Scan(
		&i.ID,
		&i.ItemID,
		&i.UserID,
		&i.QuantitySold,
		&i.SalePrice,
		&i.SaleDate,
		&i.CustomerName,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
