package teste

import (
	"backoffice/internal/domain/task_in/dto"
	"context"
	"encoding/json"
	"strconv"
	"testing"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func Test_Start(t *testing.T) {
	// 	os.Getenv("TASKS_IN_NAME")
	// 	=
	// TASKS_IN_URL=amqp://user:password@localhost:5672
	StressPublish(context.Background(), "amqp://user:password@localhost:5672", "tasks.in", 1, 1)
}

func StressPublish(ctx context.Context, url, qName string, rps, seconds int) error {
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

	task := dto.Task{}
	for i := range 100 {
		task.TaskID = strconv.Itoa(i)
		body, err := json.Marshal(&task)
		if err != nil {
			return err
		}
		err = ch.PublishWithContext(ctx, "", qName, false, false, amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
			MessageId:    strconv.FormatInt(time.Now().UnixNano(), 10),
		})
		if err != nil {
			return err
		}
	}

	return nil

}
