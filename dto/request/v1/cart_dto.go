package v1

import (
	"errors"
	"order-service/apperrors"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type AddToCartDTO struct {
	ProductID string `json:"productId" validate:"required"`
	Quantity  int    `json:"quantity" validate:"gte=1"`
}

func (dto AddToCartDTO) Validate(ctx echo.Context) error {
	validate := validator.New()
	if err := validate.Struct(dto); err != nil {
		vErr := apperrors.TryTranslateValidationErrors(err)
		return errors.New(vErr)
	}

	return nil
}
