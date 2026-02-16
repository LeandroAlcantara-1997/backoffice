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

var payloads = [][]byte{
	[]byte(`{
	"task_id": "123",
	"payload": "processar isso",
	"processing_time_ms": 500
}`),
	[]byte(`{
	teste
}`),
	[]byte(`{
	"task_id": "123",
	"payload": "processar isso",
	"processing_time_ms": 2000
}`),
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

	for i := range 10000 {
		var payload []byte
		if i%2 == 0 {
			payload, err = json.Marshal(&dto.Task{
				TaskID:           strconv.Itoa(i),
				Payload:          "processar isso",
				ProcessingTimeMS: 500,
			})
		}

		if i%3 == 0 {
			payload, err = json.Marshal(&dto.Task{
				TaskID:           strconv.Itoa(i),
				Payload:          "processar isso",
				ProcessingTimeMS: 200000,
			})
		}

		err = ch.PublishWithContext(ctx, "", qName, false, false, amqp.Publishing{
			ContentType:  "application/json",
			Body:         payload,
			DeliveryMode: amqp.Persistent,
			MessageId:    strconv.FormatInt(time.Now().UnixNano(), 10),
		})
		if err != nil {
			return err
		}
	}

	return nil
}
