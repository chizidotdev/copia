package adapters

import (
	"context"
	"github.com/chizidotdev/copia/internal/app/core"
	"github.com/chizidotdev/copia/internal/app/usecases"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"log"
	"time"
)

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
	OrderID     uuid.UUID `gorm:"not null" json:"order_id"`
	ProductName string    `gorm:"not null" json:"product_name"`
	Quantity    int64     `gorm:"not null" json:"quantity"`
	UnitPrice   float32   `gorm:"not null" json:"unit_price"`
	SubTotal    float32   `gorm:"not null" json:"sub_total"`
}

var _ usecases.OrderRepository = (*OrderRepositoryImpl)(nil)

type OrderRepositoryImpl struct {
	DB *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepositoryImpl {
	err := db.AutoMigrate(&Order{}, &OrderItem{})
	if err != nil {
		log.Panic("Failed to migrate Order and OrderItem", err)
	}
	return &OrderRepositoryImpl{DB: db}
}

func (r *OrderRepositoryImpl) ListOrders(_ context.Context, userID uuid.UUID) ([]core.Order, error) {
	var orders []core.Order
	result := r.DB.Preload("OrderItems").Find(&orders, "user_id = ?", userID)
	return orders, result.Error
}

func (r *OrderRepositoryImpl) GetOrder(_ context.Context, id uuid.UUID) (core.Order, error) {
	var order core.Order
	result := r.DB.Preload("OrderItems").First(&order, "id = ?", id)
	return order, result.Error
}

func (r *OrderRepositoryImpl) CreateOrder(_ context.Context, arg core.Order) (core.Order, error) {
	order := Order{
		Status:                arg.Status,
		EstimatedDeliveryDate: arg.EstimatedDeliveryDate,
		OrderDate:             arg.OrderDate,
		TotalAmount:           arg.TotalAmount,
		CustomerID:            arg.CustomerID,
		UserID:                arg.UserID,
	}
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&order).Error; err != nil {
			return err
		}

		orderItems := make([]OrderItem, len(arg.OrderItems))
		for i, orderItem := range arg.OrderItems {
			orderItems[i] = OrderItem{
				OrderID:     order.ID,
				ProductName: orderItem.ProductName,
				Quantity:    orderItem.Quantity,
				UnitPrice:   orderItem.UnitPrice,
				SubTotal:    orderItem.SubTotal,
			}
		}

		if err := tx.Create(orderItems).Error; err != nil {
			return err
		}

		return nil
	})
	return core.Order{
		ID:                    order.ID,
		Status:                order.Status,
		EstimatedDeliveryDate: order.EstimatedDeliveryDate,
		OrderDate:             order.OrderDate,
		TotalAmount:           order.TotalAmount,
		CustomerID:            order.CustomerID,
		UserID:                order.UserID,
	}, err
}

func (r *OrderRepositoryImpl) DeleteOrder(_ context.Context, arg core.DeleteOrderRequest) error {
	result := r.DB.Delete(&Order{}, "id = ? AND user_id = ?", arg.ID, arg.UserID)
	return result.Error
}
