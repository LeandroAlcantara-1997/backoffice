package broker

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Broker interface {
	Consumer(ctx context.Context) (<-chan amqp.Delivery, error)
	Publish(ctx context.Context, payload []byte) error
	Close() error
}

type rabbitMQ struct {
	queueName string
	channel   *amqp.Channel
	conn      *amqp.Connection
}

func NewBroker(queueName string, ch *amqp.Channel, conn *amqp.Connection) Broker {
	return &rabbitMQ{
		queueName: queueName,
		channel:   ch,
		conn:      conn,
	}
}

func (r *rabbitMQ) Consumer(ctx context.Context) (<-chan amqp.Delivery, error) {
	// QoS é importante para backpressure; ajuste conforme sua concorrência
	if err := r.channel.Qos(10, 0, false); err != nil {
		return nil, err
	}

	deliveries, err := r.channel.Consume(
		r.queueName,
		"",    // consumerTag (deixe o broker gerar)
		false, // autoAck = false  (ESSENCIAL pra não perder)
		false, // exclusive
		false, // noLocal (ignorado pelo RabbitMQ)
		false, // noWait
		nil,
	)
	if err != nil {
		return nil, err
	}

	return deliveries, nil
	// out := make(chan Delivery, 100) // buffer para absorver picos

	// go func() {
	// 	defer close(out)

	// 	for {
	// 		select {
	// 		case <-ctx.Done():
	// 			// Parar novas entregas e drenar com sua política
	// 			_ = r.channel.Cancel("", false)
	// 			for d := range deliveries {
	// 				_ = d.Nack(false, true) // requeue no shutdown (ou DLX)
	// 			}
	// 			return

	// 		case d, ok := <-deliveries:
	// 			if !ok {
	// 				// consumer cancelado / channel AMQP fechado
	// 				return
	// 			}

	// 			// Tentar repassar para 'out' respeitando o contexto.
	// 			// Se o ctx cancelar nesse momento, NÃO deixe a msg cair: Nack(requeue).
	// 			select {
	// 			case out <- Delivery{d: d}:
	// 				// entregue com sucesso para downstream; downstream fará Ack/Nack
	// 			case <-ctx.Done():
	// 				// NÃO PERDER 'd'!
	// 				_ = d.Nack(false, true)
	// 				return
	// 			}
	// 		}
	// 	}
	// }()

	// return out, nil
}

func (r *rabbitMQ) Close() error {
	if err := r.channel.Close(); err != nil {
		return err
	}

	if err := r.conn.Close(); err != nil {
		return err
	}

	return nil
}

func (r *rabbitMQ) Publish(ctx context.Context, payload []byte) error {
	return r.channel.PublishWithContext(ctx, "", r.queueName, true, false, amqp.Publishing{
		Body: payload,
	})
}
