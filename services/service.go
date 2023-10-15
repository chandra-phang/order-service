package services

import "order-service/handlers"

func InitServices(h handlers.Handler) {
	InitOrderService(h)
	InitCartService(h)
}
