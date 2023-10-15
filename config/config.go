package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var config *Config

type Config struct {
	ProductSvcHost         string
	GetProductUri          string
	IncreaseBookedQuotaUri string
	DecreaseBookedQuotaUri string

	AuthSvcHost     string
	AuthenticateUri string
}

func InitConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config = &Config{
		ProductSvcHost:         os.Getenv("PRODUCT_SERVICE_HOST"),
		GetProductUri:          os.Getenv("GET_PRODUCT_URI"),
		IncreaseBookedQuotaUri: os.Getenv("INCREASE_BOOKED_QUOTA"),
		DecreaseBookedQuotaUri: os.Getenv("DECREASE_BOOKED_QUOTA"),

		AuthSvcHost:     os.Getenv("AUTH_SERVICE_HOST"),
		AuthenticateUri: os.Getenv("AUTHENTICATE_URI"),
	}
}

func GetConfig() *Config {
	return config
}
