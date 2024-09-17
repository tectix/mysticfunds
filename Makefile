# Detect the operating system
ifeq ($(OS),Windows_NT)
    SHELL := cmd.exe
    RM := del /Q
    MKDIR := mkdir
    PSQL := psql
    MIGRATE := migrate.exe
    PROTOC := protoc.exe
else
    SHELL := /bin/bash
    RM := rm -f
    MKDIR := mkdir -p
    PSQL := psql
    MIGRATE := migrate
    PROTOC := protoc
endif

# PostgreSQL connection details
PG_HOST := localhost
PG_PORT := 5432
PG_USER := postgres
PG_PASSWORD := password

SERVICES := auth wizard mana spell realm

.PHONY: all create-dbs nuke init-migrations migrate-up migrate-down migration-status proto build run test clean help

all: create-dbs migrate-up build

create-dbs:
	@echo "Creating databases..."
	@$(foreach service,$(SERVICES),\
		PGPASSWORD=$(PG_PASSWORD) $(PSQL) -h $(PG_HOST) -p $(PG_PORT) -U $(PG_USER) -d postgres -c "CREATE DATABASE $(service);" || true;)

nuke:
	@echo "Dropping databases..."
	@$(foreach service,$(SERVICES),\
		PGPASSWORD=$(PG_PASSWORD) $(PSQL) -h $(PG_HOST) -p $(PG_PORT) -U $(PG_USER) -d postgres -c "SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = '$(service)';" && \
		PGPASSWORD=$(PG_PASSWORD) $(PSQL) -h $(PG_HOST) -p $(PG_PORT) -U $(PG_USER) -d postgres -c "DROP DATABASE IF EXISTS $(service);" || true;)

init-migrations:
	@echo "Initializing migrations..."
	@$(foreach service,$(SERVICES),\
		$(MKDIR) migrations/$(service) && \
		$(MIGRATE) create -ext sql -dir migrations/$(service) -seq init;)

migrate-up:
	@echo "Running migrations up..."
	@$(foreach service,$(SERVICES),\
		$(MIGRATE) -path migrations/$(service) -database "postgresql://$(PG_USER):$(PG_PASSWORD)@$(PG_HOST):$(PG_PORT)/$(service)?sslmode=disable" up;)

migrate-down:
	@echo "Running migrations down..."
	@$(foreach service,$(SERVICES),\
		$(MIGRATE) -path migrations/$(service) -database "postgresql://$(PG_USER):$(PG_PASSWORD)@$(PG_HOST):$(PG_PORT)/$(service)?sslmode=disable" down;)

migration-status:
	@echo "Checking migration status..."
	@$(foreach service,$(SERVICES),\
		$(MIGRATE) -path migrations/$(service) -database "postgresql://$(PG_USER):$(PG_PASSWORD)@$(PG_HOST):$(PG_PORT)/$(service)?sslmode=disable" version;)

proto:
	@echo "Generating protobuf code..."
	@$(foreach service,$(SERVICES),\
		$(PROTOC) --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/$(service)/$(service).proto;)

build:
	@echo "Building services..."
	@$(foreach service,$(SERVICES),\
		go build -o bin/$(service) cmd/$(service)/main.go;)

run:
	@echo "Running services..."
	@$(foreach service,$(SERVICES),\
		./bin/$(service) &)

test:
	@echo "Running tests..."
	@go test ./...

clean:
	@echo "Cleaning up..."
	@$(RM) bin/*

help:
	@echo "Available commands:"
	@echo "  make create-dbs    - Create databases for all services"
	@echo "  make migrate-up    - Run all migrations up"
	@echo "  make migrate-down  - Run all migrations down"
	@echo "  make migration-status - Check migration status"
	@echo "  make proto         - Generate protobuf code"
	@echo "  make build         - Build all services"
	@echo "  make run           - Run all services"
	@echo "  make test          - Run all tests"
	@echo "  make clean         - Remove built binaries"
	@echo "  make all           - Create DBs, run migrations, and build (default)"
	@echo "  make help          - Show this help message"