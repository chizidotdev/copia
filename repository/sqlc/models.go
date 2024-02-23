// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package repository

import (
	"database/sql/driver"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

type UserRole string

const (
	UserRoleMaster   UserRole = "master"
	UserRoleVendor   UserRole = "vendor"
	UserRoleCustomer UserRole = "customer"
)

func (e *UserRole) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = UserRole(s)
	case string:
		*e = UserRole(s)
	default:
		return fmt.Errorf("unsupported scan type for UserRole: %T", src)
	}
	return nil
}

type NullUserRole struct {
	UserRole UserRole `json:"user_role"`
	Valid    bool     `json:"valid"` // Valid is true if UserRole is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullUserRole) Scan(value interface{}) error {
	if value == nil {
		ns.UserRole, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.UserRole.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullUserRole) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.UserRole), nil
}

type Commission struct {
	ID               pgtype.UUID    `json:"id"`
	OrderID          pgtype.UUID    `json:"order_id"`
	UserID           pgtype.UUID    `json:"user_id"`
	CommissionAmount pgtype.Numeric `json:"commission_amount"`
	PaidStatus       string         `json:"paid_status"`
}

type Customer struct {
	ID        pgtype.UUID `json:"id"`
	StoreID   pgtype.UUID `json:"store_id"`
	FirstName string      `json:"first_name"`
	LastName  string      `json:"last_name"`
	Email     string      `json:"email"`
	Phone     string      `json:"phone"`
	Address   string      `json:"address"`
}

type Link struct {
	ID         pgtype.UUID `json:"id"`
	UserID     pgtype.UUID `json:"user_id"`
	UniqueLink string      `json:"unique_link"`
	LinkType   string      `json:"link_type"`
}

type Order struct {
	ID              pgtype.UUID      `json:"id"`
	UserID          pgtype.UUID      `json:"user_id"`
	OrderDate       pgtype.Timestamp `json:"order_date"`
	TotalAmount     pgtype.Numeric   `json:"total_amount"`
	Status          string           `json:"status"`
	PaymentStatus   string           `json:"payment_status"`
	ShippingAddress string           `json:"shipping_address"`
	CreatedAt       pgtype.Timestamp `json:"created_at"`
	UpdatedAt       pgtype.Timestamp `json:"updated_at"`
}

type OrderItem struct {
	ID        pgtype.UUID    `json:"id"`
	OrderID   pgtype.UUID    `json:"order_id"`
	ProductID pgtype.UUID    `json:"product_id"`
	Quantity  int32          `json:"quantity"`
	UnitPrice pgtype.Numeric `json:"unit_price"`
	Subtotal  pgtype.Numeric `json:"subtotal"`
}

type Product struct {
	ID            pgtype.UUID      `json:"id"`
	StoreID       pgtype.UUID      `json:"store_id"`
	Sku           string           `json:"sku"`
	Name          string           `json:"name"`
	Description   string           `json:"description"`
	Price         pgtype.Numeric   `json:"price"`
	StockQuantity int32            `json:"stock_quantity"`
	CreatedAt     pgtype.Timestamp `json:"created_at"`
	UpdatedAt     pgtype.Timestamp `json:"updated_at"`
}

type Store struct {
	ID          pgtype.UUID      `json:"id"`
	UserID      pgtype.UUID      `json:"user_id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	CreatedAt   pgtype.Timestamp `json:"created_at"`
	UpdatedAt   pgtype.Timestamp `json:"updated_at"`
}

type User struct {
	ID        pgtype.UUID      `json:"id"`
	Email     string           `json:"email"`
	FirstName string           `json:"first_name"`
	LastName  string           `json:"last_name"`
	Image     string           `json:"image"`
	Password  string           `json:"password"`
	GoogleID  pgtype.Text      `json:"google_id"`
	Role      UserRole         `json:"role"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
	UpdatedAt pgtype.Timestamp `json:"updated_at"`
}
