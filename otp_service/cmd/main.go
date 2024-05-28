package main

import (
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"github.com/m-umarr/Go_auth_service/auth_service/db"
	"github.com/m-umarr/Go_auth_service/otp_service/internal/handler"
	"github.com/m-umarr/Go_auth_service/otp_service/proto"
	"github.com/m-umarr/Go_auth_service/otp_service/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Get Twilio credentials from environment variables
	accountSid := os.Getenv("TWILIO_ACCOUNT_SID")
	authToken := os.Getenv("TWILIO_AUTH_TOKEN")
	fromPhone := os.Getenv("TWILIO_FROM_PHONE")

	// Connect to the database
	dbConn, err := db.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize the OTP service
	otpService := service.NewOTPService(accountSid, authToken, fromPhone, dbConn)

	// Create the OTP service handler
	otpHandler := handler.NewOTPHandler(otpService)

	// Create a new gRPC server
	s := grpc.NewServer()

	// Register the OTP service handler with the gRPC server
	proto.RegisterOTPServiceServer(s, otpHandler)

	// Enable reflection for debugging
	reflection.Register(s)

	// Start listening on port 50052
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen on port 50052: %v", err)
	}

	log.Printf("OTP service is running on port 50052...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
