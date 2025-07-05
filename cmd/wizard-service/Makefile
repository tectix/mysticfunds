# Wizard Service Makefile

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
SERVICE_NAME := wizard
BINARY_NAME := $(SERVICE_NAME)$(EXECUTABLE_EXTENSION)
MAIN_FILE := main.go
CONFIG_FILE := config.yaml

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

# Initialize migrations
init-migrations:
	@echo Initializing migrations for $(SERVICE_NAME) service...
	@if not exist "migrations" mkdir "migrations"
	@$(MIGRATE) create -ext sql -dir migrations -seq init_$(SERVICE_NAME)_schema

# Run database migrations
migrate-up:
	@echo Running database migrations for $(SERVICE_NAME) service...
	@$(SET_ENV) $(MIGRATE) -path migrations -database "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" up

# Rollback the last database migration
migrate-down:
	@echo Rolling back the last database migration for $(SERVICE_NAME) service...
	@$(SET_ENV) $(MIGRATE) -path migrations -database "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" down

# Check migration status
migration-status:
	@echo Checking migration status for $(SERVICE_NAME) service...
	@$(SET_ENV) $(MIGRATE) -path migrations -database "postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" version

# Run tests
test:
	@echo Running tests for $(SERVICE_NAME) service...
	@go test ..\..\internal\wizard -v

# Clean up binary
clean:
	@echo Cleaning up...
	@if exist "$(BINARY_NAME)" del /Q "$(BINARY_NAME)"

# Generate proto files
proto:
	@echo Generating proto files for $(SERVICE_NAME) service...
	@protoc --go_out=. --go_opt=paths=source_relative \
			--go-grpc_out=. --go-grpc_opt=paths=source_relative \
			proto\$(SERVICE_NAME).proto

.PHONY: build run init-migrations migrate-up migrate-down migration-status test clean proto