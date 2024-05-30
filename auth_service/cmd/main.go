package main

import (
	"log"
	"net"

	"github.com/m-umarr/Go_auth_service/auth_service/db"
	"github.com/m-umarr/Go_auth_service/auth_service/internal/client"
	"github.com/m-umarr/Go_auth_service/auth_service/internal/handler"
	"github.com/m-umarr/Go_auth_service/auth_service/internal/repository"
	"github.com/m-umarr/Go_auth_service/auth_service/messaging"
	"github.com/m-umarr/Go_auth_service/auth_service/proto"
	"github.com/m-umarr/Go_auth_service/config"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Load environment variables from .env file
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to the database
	dbConn, err := db.Connect(cfg.DSN)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	conn, err := grpc.Dial(cfg.OTPServiceAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to OTP service: %v", err)
	}
	defer conn.Close()

	// Initialize the message publisher
	publisher, err := messaging.NewPublisher(cfg.AmqpURL)
	if err != nil {
		log.Fatalf("Failed to initialize publisher: %v", err)
	}
	defer publisher.Close()

	userRepo := repository.NewUserRepository(dbConn)

	otpClient := client.NewGRPCOTPClient(conn)
	// Initialize the Auth service
	authService := handler.NewAuthService(userRepo, publisher, otpClient)

	// Create a new gRPC server
	s := grpc.NewServer()

	// Register the auth service handler with the gRPC server
	proto.RegisterAuthServiceServer(s, authService)

	reflection.Register(s)

	// Start listening on port 50051
	address := ":" + cfg.AuthServicePort
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to listen on port 50051: %v", err)
	}

	log.Printf("Auth service is running on port 50051...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
