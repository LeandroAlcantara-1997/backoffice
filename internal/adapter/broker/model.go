package broker

import amqp "github.com/rabbitmq/amqp091-go"

type Delivery struct {
	d amqp.Delivery // valor (cópia) é suficiente; a lib cuida por deliveryTag
}

func (m Delivery) Body() []byte { return m.d.Body }

func (m Delivery) Ack() error {
	// multiple=false: só esta mensagem
	return m.d.Ack(false)
}

func (m Delivery) Nack(requeue bool) error {
	// multiple=false
	return m.d.Nack(false, requeue)
}

func (m Delivery) Reject(requeue bool) error {
	return m.d.Reject(requeue)
}

func (m Delivery) Redelivered() bool { return m.d.Redelivered }
func (m Delivery) MessageId() string { return m.d.MessageId }
