package app

import (
	"order-service/api"
	"order-service/config"
	"order-service/db"
	"order-service/handlers"
	"order-service/services"
)

type Application struct {
}

// Returns a new instance of the application
func NewApplication() Application {
	return Application{}
}

func (a Application) InitApplication() {
	config.InitConfig()

	database := db.InitConnection()
	h := handlers.New(database)
	services.InitServices(h)

	api.InitRoutes()

	db.CloseConnection(database)
}
