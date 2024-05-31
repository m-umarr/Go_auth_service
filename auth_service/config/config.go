package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DSN               string
	AuthServicePort   string
	OTPServiceAddress string
	AmqpURL           string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	cfg := &Config{
		DSN:               getEnv("DSN", ""),
		AuthServicePort:   getEnv("AUTH_SERVICE_PORT", "50051"),
		OTPServiceAddress: getEnv("OTP_SERVICE_ADDRESS", ""),
		AmqpURL:           getEnv("AMQP_URL", ""),
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
