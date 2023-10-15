package model

import (
	"time"

	"github.com/labstack/echo/v4"
)

type OrderItem struct {
	ID        string
	OrderID   string
	ProductID string
	Quantity  int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type IOrderItemRepository interface {
	CreateOrderItem(ctx echo.Context, orderItem OrderItem) error
	GetOrderItems(ctx echo.Context, orderID string) ([]OrderItem, error)
}
