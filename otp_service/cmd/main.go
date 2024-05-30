package main

import (
	"log"
	"net"
	"sync"

	"github.com/m-umarr/Go_auth_service/config"
	"github.com/m-umarr/Go_auth_service/otp_service/internal/handler"
	"github.com/m-umarr/Go_auth_service/otp_service/messaging"
	"github.com/m-umarr/Go_auth_service/otp_service/proto"
	"github.com/m-umarr/Go_auth_service/otp_service/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Load environment variables from .env file
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize the OTP service
	otpService := service.NewTwilioOTPService(cfg.TwilioAccountSID, cfg.TwilioAuthToken, cfg.TwilioFromPhone)

	// Create the OTP service handler
	otpHandler := handler.NewOTPHandler(otpService)

	// Create a new gRPC server
	s := grpc.NewServer()

	proto.RegisterOTPServiceServer(s, otpHandler)

	reflection.Register(s)

	// Start listening on the specified port
	address := ":" + cfg.OtpServicePort
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", cfg.OtpServicePort, err)
	}

	log.Printf("OTP service is running on port %s...", cfg.OtpServicePort)

	// Initialize the RabbitMQ consumer
	consumer := messaging.InitializeConsumer(cfg, otpService)
	defer consumer.Close()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	wg.Wait()
}
