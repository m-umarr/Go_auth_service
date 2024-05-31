package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	TwilioAccountSID string
	TwilioAuthToken  string
	TwilioFromPhone  string
	OtpServicePort   string
	AmqpURL          string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	cfg := &Config{
		TwilioAccountSID: getEnv("TWILIO_ACCOUNT_SID", ""),
		TwilioAuthToken:  getEnv("TWILIO_AUTH_TOKEN", ""),
		TwilioFromPhone:  getEnv("TWILIO_FROM_PHONE", ""),
		OtpServicePort:   getEnv("OTP_SERVICE_PORT", "50052"),
		AmqpURL:          getEnv("AMQP_URL", ""),
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
