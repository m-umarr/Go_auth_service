package service

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"

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
	otp := generateOTP()
	message := fmt.Sprintf("Your OTP is %s", otp)
	_, err := s.twilioClient.SendSMS(phoneNumber, message)
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

func generateOTP() string {
	rand.Seed(time.Now().UnixNano())
	otp := rand.Intn(1000000)
	return fmt.Sprintf("%06d", otp)
}
