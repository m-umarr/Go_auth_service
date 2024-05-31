package service

import "github.com/m-umarr/Go_auth_service/auth_service/model"

type Auth interface {
	CreateUser(phoneNumber string) error
	VerifyUser(phoneNumber string) error
	GetUserProfile(phoneNumber string) (*model.User, error)
	LogEvent(phoneNumber string, event string) error
}

type Publisher interface {
	Publish(topic string, message []byte) error
}

type OTPClient interface {
	SendOTP(phoneNumber string) error
	VerifyOTP(phoneNumber, otp string) bool
}
