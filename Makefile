.PHONY: dev-api build-api test-api test-api-integration lint-api swag tidy clean docker-dev docker-down

API_DIR=./apps/api
BINARY=./apps/api/tmp/server
GOBIN=$(shell go env GOPATH)/bin
AIR=$(GOBIN)/air
SWAG=$(GOBIN)/swag

## dev-api: Run API with hot-reload (requires air)
dev-api:
	@test -f $(AIR) || go install github.com/air-verse/air@latest
	@cd $(API_DIR) && $(AIR) -c .air.toml

## build-api: Build the API binary
build-api:
	@echo "Building API..."
	@cd $(API_DIR) && go build -o ../../bin/server ./cmd/server
	@echo "Binary at bin/server"

## test-api: Run handler & unit tests (rápido, ~0.1s, sin Docker)
test-api:
	@cd $(API_DIR) && go test ./... -v

## test-api-integration: Run ALL tests incl. repository (~5 min, requiere Docker)
test-api-integration:
	@cd $(API_DIR) && go test -tags=integration ./... -v -count=1

## lint-api: Run Go vet
lint-api:
	@cd $(API_DIR) && go vet ./...

## tidy: Tidy Go modules
tidy:
	@cd $(API_DIR) && go mod tidy

## swag: Regenerate Swagger docs (installs swag CLI if needed)
swag:
	@test -f $(SWAG) || go install github.com/swaggo/swag/cmd/swag@latest
	@cd $(API_DIR) && $(SWAG) init -g cmd/server/main.go -o docs
	@echo "Swagger docs regenerated at apps/api/docs/"

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
