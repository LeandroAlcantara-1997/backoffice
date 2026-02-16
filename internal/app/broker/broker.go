package broker

import (
	"backoffice/internal/adapter/broker"
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type Handle interface {
	HandleFunc(ctx context.Context, d amqp091.Delivery) error
}

func Start(ctx context.Context, h Handle, cs, dlq broker.Broker) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigCh)

	slog.InfoContext(ctx, "starting consumer broker")

	for {
		msgs, err := cs.Consumer(ctx)
		if err != nil {
			return err
		}

		select {
		case sig := <-sigCh:
			slog.Info("signal received; starting graceful shutdown...", "signal", sig.String())
			return gracefulShutdown(ctx, cs)

		case <-ctx.Done():
			return gracefulShutdown(ctx, cs)
		case msg, ok := <-msgs:
			if !ok {
				continue
			}

			if err := h.HandleFunc(ctx, msg); err != nil {
				if publishErr := dlq.Publish(ctx, msg.Body); publishErr != nil {
					slog.ErrorContext(ctx, "publish to out broker failed", "error", publishErr)
					_ = msg.Nack(false, false)
					gracefulShutdown(ctx, cs)
				}

				_ = msg.Reject(false)
				continue
			}

			_ = msg.Ack(false)
		}
	}
}

func gracefulShutdown(parentCtx context.Context, cs broker.Broker) error {
	shutdownCtx, shutdownCancel := context.WithTimeout(parentCtx, 30*time.Second)
	defer shutdownCancel()

	slog.Info("initiating graceful shutdown of consumer...")

	// Pare novas entregas: cancele o consumidor / feche canal do consumidor
	if err := cs.Close(); err != nil {
		slog.Error("failed to cancel/close consumer", "err", err)
	}

	// Aqui, se você tem workers/goroutines, aguarde via WaitGroup.
	<-shutdownCtx.Done()
	if shutdownCtx.Err() == context.DeadlineExceeded {
		slog.Warn("shutdown timeout exceeded; forcing exit")
	}

	slog.Info("consumer shutdown complete")
	return errors.New("server is down")
}
