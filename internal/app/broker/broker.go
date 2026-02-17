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
)

type Handle interface {
	HandleFunc(ctx context.Context, d broker.Delivery) error
}

func Start(ctx context.Context, h Handle, cs, dlq broker.Broker) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// 1) Captura de sinais em goroutine dedicada
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	defer func() {
		signal.Stop(sigCh)
		close(sigCh)
	}()

	slog.InfoContext(ctx, "starting consumer broker")

	msgs, err := cs.Consumer(ctx)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return

			case msg, ok := <-msgs:
				if !ok {
					return
				}

				// Processamento da mensagem
				if err := h.HandleFunc(ctx, msg); err != nil {
					slog.InfoContext(ctx, "rejected message:", "error", err)

					// Publica na DLQ; se falhar, Nack e inicia shutdown
					if publishErr := dlq.Publish(ctx, msg.Body()); publishErr != nil {
						slog.ErrorContext(ctx, "publish to DLQ failed", "error", publishErr)
						_ = msg.Nack(false)
						// Cancela tudo para tentar fechar com segurança
						cancel()
						return
					}

					_ = msg.Nack(false)
					continue
				}
				_ = msg.Ack()
			}
		}
	}()

	select {
	case sig := <-sigCh:
		if sig != nil {
			slog.InfoContext(ctx, "signal received; starting graceful shutdown...", "signal", sig.String())
		}
		cancel()
	case <-ctx.Done():
	}

	gracefulShutdown(ctx, cs)

	slog.InfoContext(ctx, "consumer stopped gracefully")
	return nil
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
