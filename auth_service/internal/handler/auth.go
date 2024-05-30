package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/m-umarr/Go_auth_service/auth_service/internal/service"
	"github.com/m-umarr/Go_auth_service/auth_service/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthService struct {
	Auth      service.Auth
	Publisher service.Publisher
	OTPClient service.OTPClient
	proto.UnimplementedAuthServiceServer
}

func NewAuthService(auth service.Auth, publisher service.Publisher, otpClient service.OTPClient) *AuthService {
	return &AuthService{
		Auth:      auth,
		Publisher: publisher,
		OTPClient: otpClient,
	}
}

func (s *AuthService) SignupWithPhoneNumber(ctx context.Context, req *proto.SignupRequest) (*proto.SignupResponse, error) {
	err := s.Auth.CreateUser(req.PhoneNumber)
	if err != nil {
		if err.Error() == "user already exists" {
			return nil, status.Errorf(codes.AlreadyExists, "user already exists")
		}
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	message := map[string]string{"phone_number": req.PhoneNumber}
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to marshal message: %v", err)
	}

	err = s.Publisher.Publish("verification", messageBytes)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to publish message: %v", err)
	}

	return &proto.SignupResponse{Success: true}, nil
}

func (s *AuthService) VerifyPhoneNumber(ctx context.Context, req *proto.VerifyRequest) (*proto.VerifyResponse, error) {
	if !s.OTPClient.VerifyOTP(req.PhoneNumber, req.Otp) {
		return nil, status.Errorf(codes.InvalidArgument, "invalid OTP")
	}

	err := s.Auth.VerifyUser(req.PhoneNumber)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to verify user: %v", err)
	}

	return &proto.VerifyResponse{Success: true}, nil
}

func (s *AuthService) LoginWithPhoneNumber(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	_, err := s.Auth.GetUserProfile(req.PhoneNumber)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found: %v", err)
	}

	err = s.OTPClient.SendOTP(req.PhoneNumber)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to send OTP: %v", err)
	}

	return &proto.LoginResponse{Success: true}, nil
}

func (s *AuthService) ValidatePhoneNumberLogin(ctx context.Context, req *proto.ValidateRequest) (*proto.ValidateResponse, error) {
	if !s.OTPClient.VerifyOTP(req.PhoneNumber, req.Otp) {
		return nil, status.Errorf(codes.InvalidArgument, "invalid OTP")
	}

	_, err := s.Auth.GetUserProfile(req.PhoneNumber)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to retrieve user profile: %v", err)
	}

	err = s.Auth.LogEvent(req.PhoneNumber, "login")
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to log event: %v", err)
	}

	return &proto.ValidateResponse{Success: true}, nil
}

func (s *AuthService) GetProfile(ctx context.Context, req *proto.ProfileRequest) (*proto.ProfileResponse, error) {
	profile, err := s.Auth.GetUserProfile(req.PhoneNumber)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "profile not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to retrieve profile: %v", err)
	}

	return &proto.ProfileResponse{
		PhoneNumber: profile.PhoneNumber,
		ProfileData: profile.ProfileData,
	}, nil
}
