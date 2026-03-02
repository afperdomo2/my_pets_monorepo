.PHONY: dev-api build-api test-api lint-api clean

API_DIR=./apps/api
BINARY=./apps/api/tmp/server

## dev-api: Run API with hot-reload (requires air)
dev-api:
	@which air > /dev/null || go install github.com/air-verse/air@latest
	@cd $(API_DIR) && air -c .air.toml

## build-api: Build the API binary
build-api:
	@echo "Building API..."
	@cd $(API_DIR) && go build -o ../../bin/server ./cmd/server
	@echo "Binary at bin/server"

## test-api: Run Go tests
test-api:
	@cd $(API_DIR) && go test ./... -v

## lint-api: Run Go vet
lint-api:
	@cd $(API_DIR) && go vet ./...

## tidy: Tidy Go modules
tidy:
	@cd $(API_DIR) && go mod tidy

## clean: Remove built binaries
clean:
	@rm -rf bin/ apps/api/tmp/

## docker-dev: Start all services with docker-compose (dev)
docker-dev:
	docker compose up --build

## docker-down: Stop all services
docker-down:
	docker compose down

## help: Show this help
help:
	@grep -E '^## ' Makefile | sed 's/## /  /'
