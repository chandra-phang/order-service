package model

import (
	"order-service/lib"
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

func (OrderItem) Initialize(orderID string, productID string, quantity int) *OrderItem {
	return &OrderItem{
		ID:        lib.GenerateUUID(),
		OrderID:   orderID,
		ProductID: productID,
		Quantity:  quantity,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

type IOrderItemRepository interface {
	CreateOrderItem(ctx echo.Context, orderItem OrderItem) error
	GetOrderItems(ctx echo.Context, orderID string) ([]OrderItem, error)
}
