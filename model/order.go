package model

import (
	"order-service/lib"
	"time"

	"github.com/labstack/echo/v4"
)

type Order struct {
	ID        string
	UserID    string
	Status    OrderStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

type OrderStatus string

var (
	OrderStatusCompleted OrderStatus = "COMPLETED"
	OrderStatusCancelled OrderStatus = "CANCELLED"
)

func (Order) Initialize(userID string) *Order {
	return &Order{
		ID:        lib.GenerateUUID(),
		UserID:    userID,
		Status:    OrderStatusCompleted,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

type IOrderRepository interface {
	GetOrders(ctx echo.Context, userID string) ([]Order, error)
	GetOrder(ctx echo.Context, orderID string) (*Order, error)
	CreateOrder(ctx echo.Context, order Order) error
	CancelOrder(ctx echo.Context, orderID string) error
}
