package service

import (
	"encoding/json"
	"errors"

	"github.com/m-umarr/Go_auth_service/auth_service/model"
	"github.com/m-umarr/Go_auth_service/auth_service/repository"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidOTP        = errors.New("invalid OTP")
)

type authServiceImpl struct {
	authRepo  repository.UserRepository
	publisher Publisher
	otpClient OTPClient
}

type AuthService interface {
	SignupWithPhoneNumber(phoneNumber string) error
	VerifyPhoneNumber(phoneNumber, otp string) error
	LoginWithPhoneNumber(phoneNumber string) error
	ValidatePhoneNumberLogin(phoneNumber, otp string) error
	GetProfile(phoneNumber string) (*model.User, error)
}

func NewAuthService(authRepo repository.UserRepository, publisher Publisher, otpClient OTPClient) AuthService {
	return &authServiceImpl{
		authRepo:  authRepo,
		publisher: publisher,
		otpClient: otpClient,
	}
}

func (s *authServiceImpl) SignupWithPhoneNumber(phoneNumber string) error {
	err := s.authRepo.CreateUser(phoneNumber)
	if err != nil {
		if err == ErrUserAlreadyExists {
			return ErrUserAlreadyExists
		}
		return err
	}

	message := map[string]string{"phone_number": phoneNumber}
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	err = s.publisher.Publish("verification", messageBytes)
	if err != nil {
		return err
	}

	return nil
}

func (s *authServiceImpl) VerifyPhoneNumber(phoneNumber, otp string) error {
	if !s.otpClient.VerifyOTP(phoneNumber, otp) {
		return ErrInvalidOTP
	}

	err := s.authRepo.VerifyUser(phoneNumber)
	if err != nil {
		return err
	}

	return nil
}

func (s *authServiceImpl) LoginWithPhoneNumber(phoneNumber string) error {
	_, err := s.authRepo.GetUserProfile(phoneNumber)
	if err != nil {
		if err == ErrUserNotFound {
			return ErrUserNotFound
		}
		return err
	}

	err = s.otpClient.SendOTP(phoneNumber)
	if err != nil {
		return err
	}

	return nil
}

func (s *authServiceImpl) ValidatePhoneNumberLogin(phoneNumber, otp string) error {
	if !s.otpClient.VerifyOTP(phoneNumber, otp) {
		return ErrInvalidOTP
	}

	_, err := s.authRepo.GetUserProfile(phoneNumber)
	if err != nil {
		return err
	}

	err = s.authRepo.LogEvent(phoneNumber, "login")
	if err != nil {
		return err
	}

	return nil
}

func (s *authServiceImpl) GetProfile(phoneNumber string) (*model.User, error) {
	profile, err := s.authRepo.GetUserProfile(phoneNumber)
	if err != nil {
		return nil, err
	}

	return profile, nil
}
