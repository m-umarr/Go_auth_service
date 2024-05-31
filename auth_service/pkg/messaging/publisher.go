package messaging

import (
	"log"

	"github.com/streadway/amqp"
)

type Publisher interface {
	Publish(topic string, message []byte) error
}
type PublisherImpl struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewPublisher(amqpURL string) (*PublisherImpl, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &PublisherImpl{conn: conn, channel: ch}, nil
}

func (p *PublisherImpl) Publish(queueName string, body []byte) error {
	_, err := p.channel.QueueDeclare(
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

	err = p.channel.Publish(
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		return err
	}

	return nil
}

func (p *PublisherImpl) Close() {
	if err := p.channel.Close(); err != nil {
		log.Printf("Error closing channel: %v", err)
	}
	if err := p.conn.Close(); err != nil {
		log.Printf("Error closing connection: %v", err)
	}
}
