package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/m-umarr/Go_auth_service/otp_service/config"
	"github.com/m-umarr/Go_auth_service/otp_service/internal/handler"
	"github.com/m-umarr/Go_auth_service/otp_service/internal/service"
	"github.com/m-umarr/Go_auth_service/otp_service/pkg/messaging"
	"github.com/m-umarr/Go_auth_service/otp_service/proto"

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
	lis, err := net.Listen("tcp", fmt.Sprintf(":"+cfg.OtpServicePort))
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", cfg.OtpServicePort, err)
	}

	log.Printf("OTP service is running on port %s...", cfg.OtpServicePort)

	// Initialize the RabbitMQ consumer
	consumer := messaging.InitializeConsumer(cfg, otpService)
	defer consumer.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	<-stop
	log.Println("Received stop signal, gracefully shutting down...")
	s.GracefulStop()
	log.Println("Server stopped")
}
