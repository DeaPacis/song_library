package config

import (
	"github.com/rs/zerolog/log"
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
	log.Info().Msg("Loading application configuration")

	cfg := &Config{
		DBHost:         os.Getenv("POSTGRES_HOST"),
		DBPort:         os.Getenv("POSTGRES_PORT"),
		DBUser:         os.Getenv("POSTGRES_USER"),
		DBPassword:     os.Getenv("POSTGRES_PASSWORD"),
		DBName:         os.Getenv("POSTGRES_DB"),
		AppPort:        os.Getenv("APP_PORT"),
		ExternalAPIURL: os.Getenv("EXTERNAL_API_URL"),
	}

	log.Debug().
		Str("DBHost", cfg.DBHost).
		Str("DBPort", cfg.DBPort).
		Str("DBUser", cfg.DBUser).
		Str("DBName", cfg.DBName).
		Str("AppPort", cfg.AppPort).
		Str("ExternalAPIURL", cfg.ExternalAPIURL).
		Msg("Loaded configuration")

	if cfg.DBHost == "" || cfg.DBPort == "" || cfg.DBUser == "" || cfg.DBPassword == "" || cfg.DBName == "" {
		log.Fatal().Msg("Missing required database configuration")
	}
	if cfg.AppPort == "" {
		log.Warn().Msg("APP_PORT is not set, using default 8080")
		cfg.AppPort = "8080"
	}
	if cfg.ExternalAPIURL == "" {
		log.Warn().Msg("EXTERNAL_API_URL is not set, external API calls may fail")
	}

	log.Info().Msg("Configuration loaded successfully")
	return cfg
}
