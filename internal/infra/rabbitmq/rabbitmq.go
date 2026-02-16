package rabbitmq

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %v", msg, err)
	}
}

func New(connectionURL string) (*amqp.Connection, *amqp.Channel, error) {
	// Ex.: amqp://user:pass@host:5672/vhost
	conn, err := amqp.Dial(connectionURL)
	if err != nil {
		return nil, nil, err
	}
	// defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, err
	}

	// Ajuste de QoS para controlar paralelismo e backpressure TODO
	if err := ch.Qos(10, 0, false); err != nil {
		return nil, nil, err
	}
	return conn, ch, nil

	// // 3) Garantir fila (idempotente)
	// q, err := ch.QueueDeclare(
	// 	"tasks.in", // nome da fila
	// 	true,       // durable
	// 	false,      // autoDelete
	// 	false,      // exclusive
	// 	false,      // noWait
	// 	nil,        // args
	// )
	// failOnError(err, "Falha ao declarar fila")
	// log.Println("Fila pronta:", q.Name)

	// // 4) Publicar uma mensagem
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()

	// err = ch.PublishWithContext(ctx,
	// 	"",     // exchange ("" = default)
	// 	q.Name, // routing key (nome da fila no default exchange)
	// 	false,  // mandatory
	// 	false,  // immediate (ignorado pelo RabbitMQ)
	// 	amqp.Publishing{
	// 		ContentType:  "text/plain",
	// 		Body:         []byte("hello from Go!"),
	// 		DeliveryMode: amqp.Persistent, // 2 (persistente)
	// 	},
	// )
	// failOnError(err, "Falha ao publicar")

	// // 5) Consumir (exemplo simples)
	// msgs, err := ch.Consume(
	// 	q.Name, // queue
	// 	"",     // consumer tag
	// 	true,   // auto-ack (false = vamos dar Ack manual)
	// 	false,  // exclusive
	// 	false,  // no-local (não suportado)
	// 	false,  // no-wait
	// 	nil,    // args
	// )
	// failOnError(err, "Falha ao iniciar consumo")

	// go func() {
	// 	for d := range msgs {
	// 		log.Printf("Recebido: %s", d.Body)
	// 		// processa...
	// 		d.Ack(false) // confirma
	// 	}
	// }()

	// log.Println("Aguardando mensagens. CTRL+C para sair.")
	// select {}
}

// package main

// import (
//     "context"
//     "log"
//     "os"
//     "os/signal"
//     "syscall"
//     "time"

//     amqp "github.com/rabbitmq/amqp091-go"
// )

// func main() {
//     conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
//     if err != nil {
//         log.Fatalf("dial: %v", err)
//     }
//     defer conn.Close()

//     ch, err := conn.Channel()
//     if err != nil {
//         log.Fatalf("channel: %v", err)
//     }
//     defer ch.Close()

//

//     queueName := "minha-fila"
//     consumerTag := "consumer-1"

//     msgs, err := ch.Consume(
//         queueName,
//         consumerTag,
//         false, // autoAck=false para podermos dar Ack manual
//         false,
//         false,
//         false,
//         nil,
//     )
//     if err != nil {
//         log.Fatalf("consume: %v", err)
//     }

//     // Sinalização para shutdown
//     ctx, cancel := context.WithCancel(context.Background())
//     defer cancel()

//     // Captura SIGINT/SIGTERM
//     sigCh := make(chan os.Signal, 1)
//     signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

//     // Worker de processamento
//     done := make(chan struct{})
//     go func() {
//         defer close(done)
//         for {
//             select {
//             case d, ok := <-msgs:
//                 if !ok {
//                     // canal de entregas foi fechado (após Cancel e Close)
//                     return
//                 }
//                 // Processa a mensagem
//                 if err := process(d.Body); err != nil {
//                     // Nack e requeue ou dead-letter, conforme sua estratégia
//                     _ = d.Nack(false, true)
//                     continue
//                 }
//                 _ = d.Ack(false)
//             case <-ctx.Done():
//                 return
//             }
//         }
//     }()

//     // Espera sinal e inicia shutdown
//     <-sigCh
//     log.Println("Iniciando graceful shutdown do consumidor...")

//     // 1) Cancela o consumidor (para parar de receber novas msgs)
//     if err := ch.Cancel(consumerTag, false); err != nil {
//         log.Printf("cancel consumer: %v", err)
//     }

//     // 2) Aguarda o worker drenar e finalizar processamento atual (com timeout)
//     shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
//     defer shutdownCancel()
//     select {
//     case <-done:
//         log.Println("Consumidor finalizado com sucesso.")
//     case <-shutdownCtx.Done():
//         log.Println("Timeout ao aguardar finalização. Forçando fechamento.")
//     }

//     // 3) Fecha canal e conexão
//     _ = ch.Close()
//     _ = conn.Close()
// }

// func process(body []byte) error {
//     // TODO: sua lógica de negócio
//     time.Sleep(200 * time.Millisecond)
//     return nil
// }
