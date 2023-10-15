package api

import (
	"log"
	v1 "order-service/api/controllers/v1"
	"order-service/api/middleware"

	"github.com/labstack/echo/v4"
)

func InitRoutes() {
	e := echo.New()

	orderController := v1.InitOrderController()
	cartController := v1.InitCartController()

	v1Api := e.Group("v1")
	v1Api.Use(middleware.AuthMiddleware)
	v1Api.POST("/carts", cartController.AddToCart)
	v1Api.POST("/orders", orderController.CreateOrder)
	v1Api.PUT("/orders/:id/cancel", orderController.CancelOrder)
	v1Api.GET("/orders", orderController.ListOrder)

	log.Println("Server is running at 8082 port.")
	e.Logger.Fatal(e.Start(":8082"))
}
