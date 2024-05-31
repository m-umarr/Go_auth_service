package main

import (
	"fmt"
	"log"
	"net"

	"github.com/m-umarr/Go_auth_service/auth_service/config"
	"github.com/m-umarr/Go_auth_service/auth_service/internal/db"
	"github.com/m-umarr/Go_auth_service/auth_service/internal/handler"
	"github.com/m-umarr/Go_auth_service/auth_service/internal/repository"
	"github.com/m-umarr/Go_auth_service/auth_service/pkg/client"
	"github.com/m-umarr/Go_auth_service/auth_service/pkg/messaging"
	"github.com/m-umarr/Go_auth_service/auth_service/proto"

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

	// Initialize the message publisher

	userRepo := repository.NewUserRepository(dbConn)

	conn, err := grpc.Dial(cfg.OTPServiceAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to OTP service: %v", err)
	}
	defer conn.Close()

	otpClient := client.NewGRPCOTPClient(conn)

	publisher, err := messaging.NewPublisher(cfg.AmqpURL)
	if err != nil {
		log.Fatalf("Failed to initialize publisher: %v", err)
	}
	defer publisher.Close()

	// Initialize the Auth service
	authService := handler.NewAuthService(userRepo, publisher, otpClient)

	// Create a new gRPC server
	s := grpc.NewServer()

	// Register the auth service handler with the gRPC server
	proto.RegisterAuthServiceServer(s, authService)

	reflection.Register(s)

	// Start listening on port 50051
	lis, err := net.Listen("tcp", fmt.Sprintf(":"+cfg.AuthServicePort))
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", cfg.AuthServicePort, err)
	}

	log.Printf("Auth service is running on port %s", cfg.AuthServicePort)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
