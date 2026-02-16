.PHONY: up
up:
	@docker-compose -f ./local/docker-compose.yaml up

.PHONY: run
run:
	@go run ./cmd/main.go