package model

import (
	"time"

	"github.com/labstack/echo/v4"
)

type CartProduct struct {
	ID        string
	UserID    string
	ProductID string
	Quantity  int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ICartProductRepository interface {
	GetCartProduct(ctx echo.Context, userID string, productID string) (*CartProduct, error)
	CreateCartProduct(ctx echo.Context, cart CartProduct) error
	UpdateQuantity(ctx echo.Context, cartProductID string, quantity int) error
}
