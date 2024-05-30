package messaging

import (
	"log"

	"github.com/streadway/amqp"
)

type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
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
