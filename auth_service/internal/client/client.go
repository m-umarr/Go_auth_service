package client

import (
	"context"

	"github.com/m-umarr/Go_auth_service/otp_service/proto"
	"google.golang.org/grpc"
)

type GRPCOTPClient struct {
	client proto.OTPServiceClient
}

func NewGRPCOTPClient(conn *grpc.ClientConn) *GRPCOTPClient {
	client := proto.NewOTPServiceClient(conn)
	return &GRPCOTPClient{client: client}
}

func (c *GRPCOTPClient) SendOTP(phoneNumber string) error {
	req := &proto.GenerateOTPRequest{PhoneNumber: phoneNumber}
	_, err := c.client.GenerateOTP(context.Background(), req)
	return err
}

func (c *GRPCOTPClient) VerifyOTP(phoneNumber, otp string) bool {
	req := &proto.VerifyOTPRequest{PhoneNumber: phoneNumber, Otp: otp}
	resp, err := c.client.VerifyOTP(context.Background(), req)
	if err != nil {
		return false
	}
	return resp.Valid
}
