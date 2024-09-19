# Detect the operating system
ifeq ($(OS),Windows_NT)
    SHELL := cmd.exe
    RM := del /Q
    MKDIR := mkdir
    PSQL := psql
    MIGRATE := migrate.exe
    PROTOC := protoc.exe
    SET_ENV := set "PGPASSWORD=$(PG_PASSWORD)" &
    RUN := start /B
    NULL := nul
    SEP := &
else
    SHELL := /bin/bash
    RM := rm -f
    MKDIR := mkdir -p
    PSQL := psql
    MIGRATE := migrate
    PROTOC := protoc
    SET_ENV := export PGPASSWORD="$(PG_PASSWORD)" &&
    RUN := nohup
    NULL := /dev/null
    SEP := ;
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
		echo "Creating database for $(service)..." $(SEP) \
		$(SET_ENV) $(PSQL) -h $(PG_HOST) -p $(PG_PORT) -U $(PG_USER) -d postgres -c "CREATE DATABASE $(service);" 2>$(NULL) || echo "Database $(service) may already exist" $(SEP))

nuke:
	@echo "Dropping databases..."
	@$(foreach service,$(SERVICES),\
		echo "Dropping database $(service)..." $(SEP) \
		$(SET_ENV) $(PSQL) -h $(PG_HOST) -p $(PG_PORT) -U $(PG_USER) -d postgres -c "SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = '$(service)';" $(SEP) \
		$(PSQL) -h $(PG_HOST) -p $(PG_PORT) -U $(PG_USER) -d postgres -c "DROP DATABASE IF EXISTS $(service);" 2>$(NULL) || echo "Failed to drop $(service)" $(SEP))

init-migrations:
	@echo "Initializing migrations..."
	@$(foreach service,$(SERVICES),\
		echo "Initializing migrations for $(service)..." $(SEP) \
		$(MKDIR) migrations$(if $(findstring cmd.exe,$(SHELL)),\,/)$(service) $(SEP) \
		$(MIGRATE) create -ext sql -dir migrations/$(service) -seq init $(SEP))

migrate-up:
	@echo "Running migrations up..."
	@$(foreach service,$(SERVICES),\
		echo "Migrating $(service) up..." $(SEP) \
		$(SET_ENV) $(MIGRATE) -path migrations/$(service) -database "postgresql://$(PG_USER):$(PG_PASSWORD)@$(PG_HOST):$(PG_PORT)/$(service)?sslmode=disable" up || echo "Migration for $(service) failed" $(SEP))

migrate-down:
	@echo "Running migrations down..."
	@$(foreach service,$(SERVICES),\
		echo "Migrating $(service) down..." $(SEP) \
		$(SET_ENV) $(MIGRATE) -path migrations/$(service) -database "postgresql://$(PG_USER):$(PG_PASSWORD)@$(PG_HOST):$(PG_PORT)/$(service)?sslmode=disable" down || echo "Migration for $(service) failed" $(SEP))

migration-status:
	@echo "Checking migration status..."
	@$(foreach service,$(SERVICES),\
		echo "Checking status for $(service)..." $(SEP) \
		$(SET_ENV) $(MIGRATE) -path migrations/$(service) -database "postgresql://$(PG_USER):$(PG_PASSWORD)@$(PG_HOST):$(PG_PORT)/$(service)?sslmode=disable" version || echo "Failed to get status for $(service)" $(SEP))

proto:
	@echo "Generating protobuf code..."
	@$(foreach service,$(SERVICES),\
		echo "Generating protobuf for $(service)..." $(SEP) \
		$(PROTOC) --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/$(service)/$(service).proto $(SEP))

build:
	@echo "Building services..."
	@$(foreach service,$(SERVICES),\
		echo "Building $(service)..." $(SEP) \
		go build -o bin/$(service) cmd/$(service)/main.go $(SEP))

run:
	@echo "Running services..."
ifeq ($(OS),Windows_NT)
	@for %%s in ($(SERVICES)) do (echo Starting %%s & $(RUN) bin\%%s)
else
	@$(foreach service,$(SERVICES),\
		echo "Starting $(service)..." $(SEP) \
		$(RUN) ./bin/$(service) > $(NULL) 2>&1 & $(SEP))
endif

test:
	@echo "Running tests..."
	@go test ./...

clean:
	@echo "Cleaning up..."
	@$(RM) bin$(if $(findstring cmd.exe,$(SHELL)),\,/)* 2>$(NULL)

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