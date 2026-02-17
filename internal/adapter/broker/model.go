package broker

import amqp "github.com/rabbitmq/amqp091-go"

type Delivery struct {
	D amqp.Delivery
}

func (m Delivery) Body() []byte { return m.D.Body }

func (m Delivery) Ack() error {
	return m.D.Ack(false)
}

func (m Delivery) Nack(requeue bool) error {
	return m.D.Nack(false, requeue)
}

func (m Delivery) Reject(requeue bool) error {
	return m.D.Reject(requeue)
}

func (m Delivery) Redelivered() bool { return m.D.Redelivered }
func (m Delivery) MessageId() string { return m.D.MessageId }
