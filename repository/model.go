package repository

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Base contains common columns for all tables.
type Base struct {
	ID        uuid.UUID      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
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
	Email     string  `gorm:"not null;uniqueIndex" json:"email"`
	Password  string  `json:"password"`
	GoogleID  string  `gorm:"unique" json:"google_id"`
	Orders    []Order `gorm:"foreignKey:UserID" json:"orders"`
}

type Order struct {
	Base
	Status                string    `gorm:"not null" json:"status"`
	EstimatedDeliveryDate time.Time `gorm:"not null" json:"estimated_delivery_date"`
	OrderDate             time.Time `gorm:"not null" json:"order_date"`
	TotalAmount           float32   `gorm:"not null" json:"total_amount"`

	OrderItems []OrderItem `gorm:"foreignKey:OrderID" json:"order_items"`
	CustomerID uuid.UUID   `json:"customer_id"`
	UserID     uuid.UUID   `gorm:"not null" json:"user_id"`
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
