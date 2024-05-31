package handler

import (
	"context"
	"database/sql"
	"errors"

	"github.com/m-umarr/Go_auth_service/auth_service/proto"
	"github.com/m-umarr/Go_auth_service/auth_service/service"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServiceHandler struct {
	authService service.AuthService
	proto.UnimplementedAuthServiceServer
}

func NewAuthServiceHandler(authService service.AuthService) *AuthServiceHandler {
	return &AuthServiceHandler{
		authService: authService,
	}
}

func (s *AuthServiceHandler) SignupWithPhoneNumber(ctx context.Context, req *proto.SignupRequest) (*proto.SignupResponse, error) {
	err := s.authService.SignupWithPhoneNumber(req.PhoneNumber)
	if err != nil {
		if err == service.ErrUserAlreadyExists {
			return nil, status.Errorf(codes.AlreadyExists, "user already exists")
		}
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	return &proto.SignupResponse{Success: true}, nil
}

func (s *AuthServiceHandler) VerifyPhoneNumber(ctx context.Context, req *proto.VerifyRequest) (*proto.VerifyResponse, error) {
	err := s.authService.VerifyPhoneNumber(req.PhoneNumber, req.Otp)
	if err != nil {
		if err == service.ErrInvalidOTP {
			return nil, status.Errorf(codes.InvalidArgument, "invalid OTP")
		}
		return nil, status.Errorf(codes.Internal, "failed to verify user: %v", err)
	}

	return &proto.VerifyResponse{Success: true}, nil
}

func (s *AuthServiceHandler) LoginWithPhoneNumber(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	err := s.authService.LoginWithPhoneNumber(req.PhoneNumber)
	if err != nil {
		if err == service.ErrUserNotFound {
			return nil, status.Errorf(codes.NotFound, "user not found: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "failed to send OTP: %v", err)
	}

	return &proto.LoginResponse{Success: true}, nil
}

func (s *AuthServiceHandler) ValidatePhoneNumberLogin(ctx context.Context, req *proto.ValidateRequest) (*proto.ValidateResponse, error) {
	err := s.authService.ValidatePhoneNumberLogin(req.PhoneNumber, req.Otp)
	if err != nil {
		if err == service.ErrInvalidOTP {
			return nil, status.Errorf(codes.InvalidArgument, "invalid OTP")
		}
		return nil, status.Errorf(codes.Internal, "failed to retrieve user profile: %v", err)
	}

	return &proto.ValidateResponse{Success: true}, nil
}

func (s *AuthServiceHandler) GetProfile(ctx context.Context, req *proto.ProfileRequest) (*proto.ProfileResponse, error) {
	profile, err := s.authService.GetProfile(req.PhoneNumber)
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
