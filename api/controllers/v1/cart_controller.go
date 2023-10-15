package v1

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	controller "order-service/api/controllers"
	v1request "order-service/dto/request/v1"
	"order-service/services"

	"github.com/labstack/echo/v4"
)

type cartController struct {
	svc services.ICartService
}

// creates a new instance of this controller with reference to CartService
func InitCartController() *cartController {
	//  initializes its "svc" field with a service instance returned by "application.GetCartService()".
	return &cartController{
		svc: services.GetCartService(),
	}
}

func (c *cartController) AddToCart(ctx echo.Context) error {
	reqBody, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		return controller.WriteError(ctx, http.StatusBadRequest, err)
	}

	dto := v1request.AddToCartDTO{}
	if err := json.Unmarshal(reqBody, &dto); err != nil {
		return controller.WriteError(ctx, http.StatusBadRequest, err)
	}

	err = dto.Validate(ctx)
	if err != nil {
		return controller.WriteError(ctx, http.StatusBadRequest, err)
	}

	err = c.svc.AddToCart(ctx, dto)
	if err != nil {
		return controller.WriteError(ctx, http.StatusInternalServerError, err)
	}

	return controller.WriteSuccess(ctx, http.StatusCreated, nil)
}
