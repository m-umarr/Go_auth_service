package handler

import (
	"context"

	"github.com/m-umarr/Go_auth_service/otp_service/internal/service"
	proto "github.com/m-umarr/Go_auth_service/otp_service/proto"
)

type OTPHandler struct {
	otpService service.OTPService
	proto.UnimplementedOTPServiceServer
}

func NewOTPHandler(otpService service.OTPService) *OTPHandler {
	return &OTPHandler{otpService: otpService}
}

func (h *OTPHandler) GenerateOTP(ctx context.Context, req *proto.GenerateOTPRequest) (*proto.GenerateOTPResponse, error) {
	err := h.otpService.SendOTP(req.PhoneNumber)
	if err != nil {
		return nil, err
	}

	return &proto.GenerateOTPResponse{Success: true}, nil
}

func (h *OTPHandler) VerifyOTP(ctx context.Context, req *proto.VerifyOTPRequest) (*proto.VerifyOTPResponse, error) {
	valid := h.otpService.VerifyOTP(req.PhoneNumber, req.Otp)
	return &proto.VerifyOTPResponse{Valid: valid}, nil
}
