.PHONY: up
up:
	@docker-compose -f ./local/docker-compose.yaml up -d

.PHONY: run
run:
	@go run ./cmd/main.go

.PHONY: integration
integration:
	@go run ./test/rabbitmq.go

.PHONY: backoffice
backoffice:
	@docker build -t backoffice .
	@docker run --env-file .env --network local_backoffice-network backoffice

.PHONY: test
test:
	@go test -v ./...
