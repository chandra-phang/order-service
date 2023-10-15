package v1

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	controller "order-service/api/controllers"
	v1req "order-service/dto/request/v1"
	v1resp "order-service/dto/response/v1"
	"order-service/services"

	"github.com/labstack/echo/v4"
)

type orderController struct {
	svc services.IOrderService
}

// creates a new instance of this controller with reference to OrderService
func InitOrderController() *orderController {
	//  initializes its "svc" field with a service instance returned by "application.GetOrderService()".
	return &orderController{
		svc: services.GetOrderService(),
	}
}

func (c *orderController) CreateOrder(ctx echo.Context) error {
	reqBody, _ := ioutil.ReadAll(ctx.Request().Body)
	dto := v1req.CreateOrderDTO{}

	if err := json.Unmarshal(reqBody, &dto); err != nil {
		return controller.WriteError(ctx, http.StatusBadRequest, err)
	}

	err := dto.Validate(ctx)
	if err != nil {
		return controller.WriteError(ctx, http.StatusBadRequest, err)
	}

	err = c.svc.CreateOrder(ctx, dto)
	if err != nil {
		return controller.WriteError(ctx, http.StatusInternalServerError, err)
	}

	return controller.WriteSuccess(ctx, http.StatusCreated, nil)
}

func (c *orderController) CancelOrder(ctx echo.Context) error {
	productID := ctx.Param("id")

	err := c.svc.CancelOrder(ctx, productID)
	if err != nil {
		return controller.WriteError(ctx, http.StatusInternalServerError, err)
	}

	return controller.WriteSuccess(ctx, http.StatusOK, nil)
}

func (c *orderController) ListOrder(ctx echo.Context) error {
	orders, err := c.svc.ListOrder(ctx)
	if err != nil {
		return controller.WriteError(ctx, http.StatusInternalServerError, err)
	}

	resp := new(v1resp.ListOrderDTO).ConvertFromOrdersEntity(orders)
	return controller.WriteSuccess(ctx, http.StatusOK, resp)
}
