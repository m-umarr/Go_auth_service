package messaging

import (
	"encoding/json"
	"log"

	"github.com/m-umarr/Go_auth_service/otp_service/config"
	"github.com/m-umarr/Go_auth_service/otp_service/service"

	"github.com/streadway/amqp"
)

type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

type SendOTPMessage struct {
	PhoneNumber string `json:"phone_number"`
}

func NewConsumer(amqpURL string) (*Consumer, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &Consumer{conn: conn, channel: ch}, nil
}

func (c *Consumer) Consume(queueName string, handler func([]byte) error) error {
	_, err := c.channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	msgs, err := c.channel.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			if err := handler(d.Body); err != nil {
				log.Printf("Error handling message: %v", err)
			}
		}
	}()

	return nil
}

func (c *Consumer) Close() {
	if err := c.channel.Close(); err != nil {
		log.Printf("Error closing channel: %v", err)
	}
	if err := c.conn.Close(); err != nil {
		log.Printf("Error closing connection: %v", err)
	}
}

func InitializeConsumer(cfg *config.Config, otpService *service.TwilioOTPService) *Consumer {
	consumer, err := NewConsumer(cfg.AmqpURL)
	if err != nil {
		log.Fatalf("Failed to initialize consumer: %v", err)
	}

	err = consumer.Consume("verification", func(body []byte) error {
		var msg SendOTPMessage
		if err := json.Unmarshal(body, &msg); err != nil {
			return err
		}

		return otpService.SendOTP(msg.PhoneNumber)
	})
	if err != nil {
		log.Fatalf("Failed to consume messages: %v", err)
	}

	return consumer
}
