package v1

import (
	"errors"
	"order-service/apperrors"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type CreateOrderDTO struct {
	OrderItems []OrderItemDTO `json:"orderItems"`
}

type OrderItemDTO struct {
	ProductID string `json:"productId" validate:"required"`
	Quantity  int    `json:"quantity" validate:"required"`
}

func (dto CreateOrderDTO) Validate(ctx echo.Context) error {
	validate := validator.New()
	if err := validate.Struct(dto); err != nil {
		vErr := apperrors.TryTranslateValidationErrors(err)
		return errors.New(vErr)
	}

	return nil
}
