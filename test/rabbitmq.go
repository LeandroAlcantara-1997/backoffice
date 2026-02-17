package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	if err := StressPublish(context.Background(), "amqp://user:password@localhost:5672", "tasks.in", 1, 1); err != nil {
		log.Fatal(err)
	}
}

func StressPublish(ctx context.Context, url, qName string, rps, seconds int) error {
	content, err := os.ReadFile("test/integration_test.json")
	if err != nil {
		return err
	}

	conn, err := amqp.Dial(url)
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	err = ch.PublishWithContext(ctx, "", qName, false, false, amqp.Publishing{
		ContentType:  "application/json",
		Body:         content,
		DeliveryMode: amqp.Persistent,
		MessageId:    strconv.FormatInt(time.Now().UnixNano(), 10),
	})
	if err != nil {
		return err
	}

	return nil
}
