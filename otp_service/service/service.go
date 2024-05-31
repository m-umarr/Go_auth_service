package service

import (
	"github.com/m-umarr/Go_auth_service/otp_service/pkg/twilio"
)

type TwilioOTPService struct {
	twilioClient *twilio.TwilioClient
}

type OTPService interface {
	SendOTP(phoneNumber string) error
	VerifyOTP(phoneNumber, otp string) bool
}

func NewTwilioOTPService(accountSid, authToken, fromPhone string) *TwilioOTPService {
	client := twilio.NewClient(accountSid, authToken, fromPhone)
	return &TwilioOTPService{twilioClient: client}
}

func (s *TwilioOTPService) SendOTP(phoneNumber string) error {
	_, err := s.twilioClient.SendSMS(phoneNumber)
	if err != nil {
		return err
	}

	return nil
}

func (s *TwilioOTPService) VerifyOTP(phoneNumber, otp string) bool {
	err := s.twilioClient.OtpVerification(phoneNumber, otp)
	if err != nil {
		return false
	}
	return true
}
