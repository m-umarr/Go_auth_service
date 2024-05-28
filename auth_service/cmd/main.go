package main

import (
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"github.com/m-umarr/Go_auth_service/auth_service/db"
	"github.com/m-umarr/Go_auth_service/auth_service/internal/handler"
	"github.com/m-umarr/Go_auth_service/auth_service/proto"
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

	// Initialize the Auth service
	authService := handler.NewAuthService(dbConn, otpService)

	// Create a new gRPC server
	s := grpc.NewServer()

	// Register the Auth service handler with the gRPC server
	proto.RegisterAuthServiceServer(s, authService)

	// Enable reflection for debugging
	reflection.Register(s)

	// Start listening on port 50051
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen on port 50051: %v", err)
	}

	log.Printf("Auth service is running on port 50051...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
