# Detect the operating system
ifeq ($(OS),Windows_NT)
    SHELL := cmd.exe
    RM := del /Q
    RMDIR := rmdir /S /Q
    MKDIR := mkdir
    EXECUTABLE_EXTENSION := .exe
    MIGRATE := migrate.exe
    SET_ENV := set "PGPASSWORD=$(DB_PASSWORD)" &
    RUN := start /B
    NULL := nul
    SEP := &
else
    SHELL := /bin/bash
    RM := rm -f
    RMDIR := rm -rf
    MKDIR := mkdir -p
    EXECUTABLE_EXTENSION :=
    MIGRATE := migrate
    SET_ENV := export PGPASSWORD="$(DB_PASSWORD)" &&
    RUN := nohup
    NULL := /dev/null
    SEP := ;
endif

# Variables
SERVICE_NAME := mana
BINARY_NAME := $(SERVICE_NAME)$(EXECUTABLE_EXTENSION)
MAIN_FILE := main.go
CONFIG_FILE := config.yaml
MIGRATION_DIR := ../../migrations/$(SERVICE_NAME)

# Go related variables
GOBASE := $(shell cd)
GOBIN := $(GOBASE)

# Database configuration
DB_HOST := localhost
DB_PORT := 5432
DB_USER := mysticfunds
DB_PASSWORD := mysticfunds
DB_NAME := $(SERVICE_NAME)

# Build the binary
build:
	@echo "Building $(SERVICE_NAME) service..."
	@$(MKDIR) bin
	@go build -o bin/$(SERVICE_NAME)$(EXECUTABLE_EXTENSION) main.go

# Run the service
run:
	@echo "Running $(SERVICE_NAME) service..."
	@if [ -f "bin/$(SERVICE_NAME)$(EXECUTABLE_EXTENSION)" ]; then \
		bin/$(SERVICE_NAME)$(EXECUTABLE_EXTENSION); \
	else \
		echo "Service $(SERVICE_NAME) not built"; \
	fi

# Initialize migrations directory
init-migrations:
	@echo Initializing migrations for $(SERVICE_NAME) service...
	@if not exist "$(MIGRATION_DIR)" $(MKDIR) "$(MIGRATION_DIR)"

# Create a new migration
new-migration:
	@echo Creating new migration...
	@if not exist "$(MIGRATION_DIR)" $(MKDIR) "$(MIGRATION_DIR)"
	@$(MIGRATE) create -ext sql -dir $(MIGRATION_DIR) -seq $(name)

# Run database migrations
migrate-up:
	@echo Running database migrations for $(SERVICE_NAME) service...
	@$(SET_ENV) $(MIGRATE) -path $(MIGRATION_DIR) -database "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" up

# Rollback the last database migration
migrate-down:
	@echo Rolling back the last database migration for $(SERVICE_NAME) service...
	@$(SET_ENV) $(MIGRATE) -path $(MIGRATION_DIR) -database "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" down

# Check migration status
migration-status:
	@echo Checking migration status for $(SERVICE_NAME) service...
	@$(SET_ENV) $(MIGRATE) -path $(MIGRATION_DIR) -database "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" version

# Run tests
test:
	@echo Running tests for $(SERVICE_NAME) service...
	@go test ..\..\internal\mana -v

# Clean up binary and generated files
clean:
	@echo Cleaning up...
	@if exist "$(BINARY_NAME)" del /Q "$(BINARY_NAME)"

# Generate proto files
proto:
	@echo Generating proto files for $(SERVICE_NAME) service...
	@protoc --go_out=. --go_opt=paths=source_relative \
			--go-grpc_out=. --go-grpc_opt=paths=source_relative \
			../../proto/$(SERVICE_NAME)/$(SERVICE_NAME).proto

# Help
help:
	@echo Available commands:
	@echo   make build              - Build the service
	@echo   make run               - Run the service
	@echo   make init-migrations   - Initialize migrations directory
	@echo   make new-migration name=migration_name - Create a new migration
	@echo   make migrate-up        - Run all pending migrations
	@echo   make migrate-down      - Rollback the last migration
	@echo   make migration-status  - Check migration status
	@echo   make test              - Run tests
	@echo   make clean             - Clean up binaries
	@echo   make proto             - Generate proto files

.PHONY: build run init-migrations migrate-up migrate-down migration-status test clean proto new-migration help