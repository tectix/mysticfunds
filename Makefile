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
    BINARY_EXT := .exe
    CREATE_DIR = if not exist "$(1)" $(MKDIR) "$(1)"
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
    BINARY_EXT :=
    CREATE_DIR = mkdir -p $(1)
endif

# PostgreSQL connection details
PG_HOST := localhost
PG_PORT := 5432
PG_USER := postgres
PG_PASSWORD := password

SERVICES := auth wizard mana spell realm

.PHONY: all create-dbs nuke init-migrations migrate-up migrate-down migration-status proto build run test clean help

all: create-dbs migrate-up build

build:
	@echo "Building services..."
	@for %%s in ($(SERVICES)) do ( \
		echo Building %%s... & \
		if not exist "cmd\%%s-service\bin" $(MKDIR) "cmd\%%s-service\bin" & \
		cd cmd/%%s-service && go build -o "bin/%%s$(BINARY_EXT)" main.go && cd ../.. \
	)

run:
	@echo "Running services..."
	@for %%s in ($(SERVICES)) do ( \
		echo Starting %%s & \
		if exist "cmd\%%s-service\bin\%%s$(BINARY_EXT)" ( \
			cd cmd/%%s-service && $(MAKE) run && cd ../.. \
		) else ( \
			echo Service %%s not built \
		) \
	)

# Create service-specific run targets
define make-run-target
run-$(1):
	@echo "Running $(1) service..."
	@cd cmd/$(1)-service && $(MAKE) run && cd ../..
endef

# Create run targets for each service
$(foreach service,$(SERVICES),$(eval $(call make-run-target,$(service))))

# Create service-specific build targets
define make-build-target
build-$(1):
	@echo "Building $(1) service..."
	@cd cmd/$(1)-service && $(MAKE) build && cd ../..
endef

# Create build targets for each service
$(foreach service,$(SERVICES),$(eval $(call make-build-target,$(service))))

create-dbs:
	@echo "Creating databases..."
	@for %%s in ($(SERVICES)) do ( \
		echo Creating database for %%s... & \
		$(SET_ENV) $(PSQL) -h $(PG_HOST) -p $(PG_PORT) -U $(PG_USER) -d postgres -c "CREATE DATABASE %%s;" 2>$(NULL) || echo Database %%s may already exist \
	)

nuke:
	@echo "Dropping databases..."
	@for %%s in ($(SERVICES)) do ( \
		echo Dropping database %%s... & \
		$(SET_ENV) $(PSQL) -h $(PG_HOST) -p $(PG_PORT) -U $(PG_USER) -d postgres -c "SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = '%%s';" & \
		$(PSQL) -h $(PG_HOST) -p $(PG_PORT) -U $(PG_USER) -d postgres -c "DROP DATABASE IF EXISTS %%s;" 2>$(NULL) || echo Failed to drop %%s \
	)

migrate-up:
	@echo "Running migrations up..."
	@for %%s in ($(SERVICES)) do ( \
		echo Migrating %%s up... & \
		$(SET_ENV) $(MIGRATE) -path migrations/%%s -database "postgresql://$(PG_USER):$(PG_PASSWORD)@$(PG_HOST):$(PG_PORT)/%%s?sslmode=disable" up || echo Migration for %%s failed \
	)

migrate-down:
	@echo "Running migrations down..."
	@for %%s in ($(SERVICES)) do ( \
		echo Migrating %%s down... & \
		$(SET_ENV) $(MIGRATE) -path migrations/%%s -database "postgresql://$(PG_USER):$(PG_PASSWORD)@$(PG_HOST):$(PG_PORT)/%%s?sslmode=disable" down || echo Migration for %%s failed \
	)

migration-status:
	@echo "Checking migration status..."
	@for %%s in ($(SERVICES)) do ( \
		echo Checking status for %%s... & \
		$(SET_ENV) $(MIGRATE) -path migrations/%%s -database "postgresql://$(PG_USER):$(PG_PASSWORD)@$(PG_HOST):$(PG_PORT)/%%s?sslmode=disable" version || echo Failed to get status for %%s \
	)

proto:
	@echo "Generating protobuf code..."
	@for %%s in ($(SERVICES)) do ( \
		echo Generating protobuf for %%s... & \
		$(PROTOC) --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/%%s/%%s.proto \
	)

test:
	@echo "Running tests..."
	@go test ./...

clean:
	@echo "Cleaning up..."
	@for %%s in ($(SERVICES)) do ( \
		echo Cleaning %%s... & \
		if exist "cmd\%%s-service\bin" rmdir /s /q "cmd\%%s-service\bin" \
	)

help:
	@echo "Available commands:"
	@echo "  make build            - Build all services"
	@echo "  make build-{service}  - Build specific service (e.g., make build-mana)"
	@echo "  make run             - Run all services"
	@echo "  make run-{service}   - Run specific service (e.g., make run-mana)"
	@echo "  make create-dbs      - Create databases for all services"
	@echo "  make migrate-up      - Run all migrations up"
	@echo "  make migrate-down    - Run all migrations down"
	@echo "  make migration-status - Check migration status"
	@echo "  make proto           - Generate protobuf code"
	@echo "  make test            - Run all tests"
	@echo "  make clean           - Clean all service binaries"
	@echo "  make all             - Create DBs, run migrations, and build (default)"
	@echo "  make help            - Show this help message"