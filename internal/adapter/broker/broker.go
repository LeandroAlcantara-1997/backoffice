package broker

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type Broker interface {
	Consumer(ctx context.Context) (<-chan Delivery, error)
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

func (r *rabbitMQ) Consumer(ctx context.Context) (<-chan Delivery, error) {
	// QoS é importante para backpressure; ajuste conforme sua concorrência
	if err := r.channel.Qos(10, 0, false); err != nil {
		return nil, err
	}

	msgs, err := r.channel.Consume(
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

	// wg := new(sync.WaitGroup)
	// wg.Add(3)
	var out = make(chan Delivery, 1024)
	go func() {
		defer close(out) // feche quando in acabar/cancelar
		for {
			select {
			case <-ctx.Done():
				return
			case v, ok := <-msgs:
				if !ok {
					return
				}
				// envio respeitando cancelamento
				select {
				case out <- Delivery{v}:
				case <-ctx.Done():
					return
				}
			}
		}
	}()
	// wg.Wait()
	// defer wg.Done()

	return out, nil
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
	ctx, span := otel.Tracer(r.queueName).Start(ctx, "publing message",
		trace.WithAttributes(attribute.String("msg", string(payload))))
	defer span.End()
	return r.channel.PublishWithContext(ctx, "", r.queueName, false, false, amqp.Publishing{
		Body: payload,
	})
}
