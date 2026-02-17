package main

import (
	"backoffice/internal/app/broker"
	tasksin "backoffice/internal/app/broker/consumer/tasks_in"
	tasksout "backoffice/internal/app/broker/consumer/tasks_out"
	"backoffice/internal/container"
	"context"
	"log/slog"
)

func main() {
	var ctx = context.Background()
	ctn, err := container.New(ctx)
	if err != nil {
		slog.Error("setup container error", "error", err)
	}
	go broker.Start(ctx,
		tasksin.NewTasksIn(ctn.TaskInUseCase),
		ctn.TaskInBroker, ctn.TaskDQLBroker)

	if err := broker.Start(ctx,
		tasksout.NewTasksOut(ctn.TaskOutUseCase),
		ctn.TaskOutBroker, ctn.TaskDQLBroker); err != nil {
		if err := ctn.CloseConnections(); err != nil {
			slog.ErrorContext(ctx, "error closing connection", "error", err)
		}
		slog.Error("start broker error", "error", err)
	}
}
