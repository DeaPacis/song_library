package config

import (
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"os"
)

type Config struct {
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	AppPort        string
	ExternalAPIURL string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Warn("No .env file found")
	}

	return &Config{
		DBHost:         os.Getenv("POSTGRES_HOST"),
		DBPort:         os.Getenv("POSTGRES_PORT"),
		DBUser:         os.Getenv("POSTGRES_USER"),
		DBPassword:     os.Getenv("POSTGRES_PASSWORD"),
		DBName:         os.Getenv("POSTGRES_DB"),
		AppPort:        os.Getenv("APP_PORT"),
		ExternalAPIURL: os.Getenv("EXTERNAL_API_URL"),
	}
}
