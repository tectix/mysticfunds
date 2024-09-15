# Wizard Service Makefile

# Variables
SERVICE_NAME := wizard-service
BINARY_NAME := $(SERVICE_NAME)
MAIN_FILE := main.go
CONFIG_FILE := config.yaml

# Go related variables
GOBASE := $(shell pwd)
GOBIN := $(GOBASE)/bin

# Database configuration
DB_HOST := localhost
DB_PORT := 5432
DB_USER := postgres
DB_PASSWORD := password
DB_NAME := wizard

# Build the binary
build:
	@echo "Building $(SERVICE_NAME)..."
	@go build -o $(GOBIN)/$(BINARY_NAME) $(MAIN_FILE)

# Run the service
run: build
	@echo "Running $(SERVICE_NAME)..."
	@$(GOBIN)/$(BINARY_NAME)

# Run database migrations
migrate:
	@echo "Running database migrations for $(SERVICE_NAME)..."
	@go run ../../cmd/migrator/main.go -config $(CONFIG_FILE) -direction up

# Rollback the last database migration
migrate-down:
	@echo "Rolling back the last database migration for $(SERVICE_NAME)..."
	@go run ../../cmd/migrator/main.go -config $(CONFIG_FILE) -direction down

# Run tests
test:
	@echo "Running tests for $(SERVICE_NAME)..."
	@go test ../../internal/wizard -v

# Clean up binary
clean:
	@echo "Cleaning up..."
	@if exist "$(GOBIN)/$(BINARY_NAME)" del /F "$(GOBIN)/$(BINARY_NAME)"
	@if exist "$(GOBIN)/$(BINARY_NAME).exe" del /F "$(GOBIN)/$(BINARY_NAME).exe"

# Generate proto files
proto:
	@echo "Generating proto files..."
	@protoc --go_out=../../ --go_opt=paths=source_relative \
			--go-grpc_out=../../ --go-grpc_opt=paths=source_relative \
			../../proto/wizard/wizard.proto

.PHONY: build run migrate migrate-down test clean proto