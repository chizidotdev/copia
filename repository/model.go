package repository

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Base contains common columns for all tables.
type Base struct {
	ID        uuid.UUID      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (base *Base) BeforeCreate(db *gorm.DB) error {
	base.ID = uuid.New()

	return nil
}

type User struct {
	Base
	FirstName string  `gorm:"not null" json:"first_name"`
	LastName  string  `gorm:"not null" json:"last_name"`
	UserEmail string  `gorm:"not null,unique" json:"email"`
	Password  string  `gorm:"not null" json:"password"`
	Orders    []Order `gorm:"foreignKey:UserID" json:"orders"`
}

type Order struct {
	Base
	Status                string    `gorm:"not null" json:"status"`
	ShippingDetails       string    `gorm:"not null" json:"shipping_details"`
	EstimatedDeliveryDate time.Time `gorm:"not null" json:"estimated_delivery_date"`
	OrderDate             time.Time `gorm:"not null" json:"order_date"`
	TotalAmount           float32   `gorm:"not null" json:"total_amount"`
	PaymentStatus         string    `gorm:"not null" json:"payment_status"`
	PaymentMethod         string    `gorm:"not null" json:"payment_method"`
	BillingAddress        string    `gorm:"not null" json:"billing_address"`
	ShippingAddress       string    `gorm:"not null" json:"shipping_address"`
	Notes                 string    `gorm:"not null" json:"notes"`

	OrderItems []OrderItem `gorm:"foreignKey:OrderID" json:"order_items"`
	CustomerID uuid.UUID   `json:"customer_id"`
	UserID     string      `gorm:"not null" json:"user_id"`
}

type OrderItem struct {
	Base
	OrderID   uuid.UUID `gorm:"not null" json:"order_id"`
	ProductID uuid.UUID `gorm:"not null" json:"product_id"`
	Quantity  int64     `gorm:"not null" json:"quantity"`
	UnitPrice float32   `gorm:"not null" json:"unit_price"`
	SubTotal  float32   `gorm:"not null" json:"sub_total"`
}

type Products struct {
	Base
	Name         string  `gorm:"not null" json:"name"`
	Description  string  `gorm:"not null" json:"description"`
	Category     string  `gorm:"not null" json:"category"`
	Price        float32 `gorm:"not null" json:"price"`
	InStock      bool    `gorm:"not null" json:"in_stock"`
	Availability string  `gorm:"not null" json:"availability"`
	ImageURL     string  `gorm:"not null" json:"image_url"`
}

type Customer struct {
	Base
	// Orders []Order `gorm:"foreignKey:CustomerID" json:"orders"`
}
