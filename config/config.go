package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var config *Config

type Config struct {
	ProductSvcHost string
	AuthSvcHost    string
}

func InitConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config = &Config{
		ProductSvcHost: os.Getenv("PRODUCT_SERVICE_HOST"),
		AuthSvcHost:    os.Getenv("AUTH_SERVICE_HOST"),
	}
}

func GetConfig() *Config {
	return config
}
