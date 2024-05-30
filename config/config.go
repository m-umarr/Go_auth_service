package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AuthServicePort   string
	OtpServicePort    string
	DSN               string
	TwilioAccountSID  string
	TwilioAuthToken   string
	TwilioFromPhone   string
	TwilioServiceID   string
	AmqpURL           string
	OTPServiceAddress string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	config := &Config{
		AuthServicePort:   getEnv("AUTH_SERVICE_PORT", "50051"),
		OtpServicePort:    getEnv("OTP_SERVICE_PORT", "50052"),
		DSN:               getEnv("DSN", ""),
		TwilioAccountSID:  getEnv("TWILIO_ACCOUNT_SID", ""),
		TwilioAuthToken:   getEnv("TWILIO_AUTH_TOKEN", ""),
		TwilioFromPhone:   getEnv("TWILIO_FROM_PHONE", ""),
		TwilioServiceID:   getEnv("TWILIO_SERVICE_ID", ""),
		AmqpURL:           getEnv("AMQP_URL", ""),
		OTPServiceAddress: getEnv("OTP_SERVICE_ADDRESS", ""),
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}
