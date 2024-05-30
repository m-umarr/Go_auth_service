package service

import (
	"database/sql"

	"github.com/m-umarr/Go_auth_service/otp_service/internal/twilio"
)

type OTPService struct {
	twilioClient *twilio.Client
}

func NewOTPService(accountSid, authToken, fromPhone string, db *sql.DB) *OTPService {
	client := twilio.NewClient(accountSid, authToken, fromPhone)
	return &OTPService{twilioClient: client}
}

func (s *OTPService) SendOTP(phoneNumber string) error {
	_, err := s.twilioClient.SendSMS(phoneNumber)
	if err != nil {
		return err
	}

	return nil
}

func (s *OTPService) VerifyOTP(phoneNumber, otp string) bool {
	err := s.twilioClient.OtpVerification(phoneNumber, otp)
	if err != nil {
		return false
	}
	return true
}
