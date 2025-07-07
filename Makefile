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
PG_USER := mysticfunds
PG_PASSWORD := mysticfunds

SERVICES := auth wizard mana
# FUTURE_SERVICES := spell realm
GATEWAY := api-gateway

.PHONY: all create-dbs nuke init-migrations migrate-up migrate-down migration-status proto build run test clean help start stop status dev logs

all: create-dbs migrate-up build

build:
	@echo "Building services..."
	@for service in $(SERVICES); do \
		echo "Building $$service..."; \
		cd cmd/$$service-service && $(MAKE) build && cd ../..; \
	done
	@echo "Building API Gateway..."
	@cd cmd/$(GATEWAY) && $(MAKE) build && cd ../..

run:
	@echo "Running services..."
	@$(call CREATE_DIR,logs)
	@echo "Starting auth service on :50051"
	@cd cmd/auth-service && nohup ./bin/auth > ../../logs/auth-service.log 2>&1 & cd ../..
	@echo "Starting wizard service on :50052"  
	@cd cmd/wizard-service && nohup ./bin/wizard > ../../logs/wizard-service.log 2>&1 & cd ../..
	@echo "Starting mana service on :50053"
	@cd cmd/mana-service && nohup ./bin/mana > ../../logs/mana-service.log 2>&1 & cd ../..
	@echo "Starting API Gateway on :8080"
	@cd cmd/api-gateway && nohup ./bin/api-gateway > ../../logs/api-gateway.log 2>&1 & cd ../..
	@sleep 3

# Create service-specific run targets
define make-run-target
run-$(1):
	@echo "Running $(1) service..."
	@cd cmd/$(1)-service && $(MAKE) run && cd ../..
endef

# Create run targets for each service
$(foreach service,$(SERVICES),$(eval $(call make-run-target,$(service))))

# Create run target for API Gateway
run-gateway:
	@echo "Running API Gateway..."
	@cd cmd/$(GATEWAY) && $(MAKE) run && cd ../..

# Create service-specific build targets
define make-build-target
build-$(1):
	@echo "Building $(1) service..."
	@cd cmd/$(1)-service && $(MAKE) build && cd ../..
endef

# Create build targets for each service
$(foreach service,$(SERVICES),$(eval $(call make-build-target,$(service))))

# Create build target for API Gateway
build-gateway:
	@echo "Building API Gateway..."
	@cd cmd/$(GATEWAY) && $(MAKE) build && cd ../..

create-dbs:
	@echo "Creating databases..."
	@for service in $(SERVICES); do \
		echo "Creating database for $$service..."; \
		$(SET_ENV) $(PSQL) -h $(PG_HOST) -p $(PG_PORT) -U $(PG_USER) -d postgres -c "CREATE DATABASE $$service;" 2>$(NULL) || echo "Database $$service may already exist"; \
	done

nuke:
	@echo "Dropping databases..."
	@for service in $(SERVICES); do \
		echo "Dropping database $$service..."; \
		$(SET_ENV) $(PSQL) -h $(PG_HOST) -p $(PG_PORT) -U $(PG_USER) -d postgres -c "SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = '$$service';" ; \
		$(PSQL) -h $(PG_HOST) -p $(PG_PORT) -U $(PG_USER) -d postgres -c "DROP DATABASE IF EXISTS $$service;" 2>$(NULL) || echo "Failed to drop $$service"; \
	done

migrate-up:
	@echo "Running migrations up..."
	@for service in $(SERVICES); do \
		echo "Migrating $$service up..."; \
		$(SET_ENV) $(MIGRATE) -path migrations/$$service -database "postgresql://$(PG_USER):$(PG_PASSWORD)@$(PG_HOST):$(PG_PORT)/$$service?sslmode=disable" up || echo "Migration for $$service failed"; \
	done

migrate-down:
	@echo "Running migrations down..."
	@for service in $(SERVICES); do \
		echo "Migrating $$service down..."; \
		$(SET_ENV) $(MIGRATE) -path migrations/$$service -database "postgresql://$(PG_USER):$(PG_PASSWORD)@$(PG_HOST):$(PG_PORT)/$$service?sslmode=disable" down || echo "Migration for $$service failed"; \
	done

migration-status:
	@echo "Checking migration status..."
	@for service in $(SERVICES); do \
		echo "Checking status for $$service..."; \
		$(SET_ENV) $(MIGRATE) -path migrations/$$service -database "postgresql://$(PG_USER):$(PG_PASSWORD)@$(PG_HOST):$(PG_PORT)/$$service?sslmode=disable" version || echo "Failed to get status for $$service"; \
	done

proto:
	@echo "Generating protobuf code..."
	@for service in $(SERVICES); do \
		echo "Generating protobuf for $$service..."; \
		$(PROTOC) --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/$$service/$$service.proto ; \
	done

test:
	@echo "Running tests..."
	@go test ./...

clean:
	@echo "Cleaning up..."
	@for service in $(SERVICES); do \
		echo "Cleaning $$service..."; \
		$(RM) -rf cmd/$$service-service/bin ; \
	done
	@echo "Cleaning API Gateway..."
	@$(RM) -rf cmd/$(GATEWAY)/bin

# Quick start shortcuts
start: stop all
	@echo "Starting MysticFunds system..."
	@echo "=============================="
	@echo ""
	@echo "  █▀▀ ▀█▀ ▄▀█ █▀█ ▀█▀ █ █▄█ █▀▀"
	@echo "  ▀▀█  █  █▀█ █▀▄  █  █ █▀█ █▄█"
	@echo ""
	@echo "=============================="
	@echo "Starting all services..."
	@$(call CREATE_DIR,logs)
	@$(MAKE) -s run
	@echo ""
	@echo "MysticFunds is running!"
	@echo "Web Interface: http://localhost:8080"
	@echo "Auth Service:  localhost:50051"
	@echo "Wizard Service: localhost:50052"
	@echo "Mana Service:  localhost:50053"
	@echo ""

stop:
	@echo "Stopping MysticFunds system..."
	@echo "=============================="
	@echo ""
	@echo "  █▀▀ ▀█▀ █▀█ █▀█ █▀█ █ █▄█ █▀▀"
	@echo "  ▀▀█  █  █▄█ █▀▀ █▀▀ █ █▀█ █▄█"
	@echo ""
	@echo "=============================="
	@pkill -f "bin/auth" || true
	@pkill -f "bin/wizard" || true
	@pkill -f "bin/mana" || true
	@pkill -f "bin/api-gateway" || true
	@sleep 2
	@pkill -9 -f "bin/auth" || true
	@pkill -9 -f "bin/wizard" || true
	@pkill -9 -f "bin/mana" || true
	@pkill -9 -f "bin/api-gateway" || true
	@echo "All services stopped"

status:
	@echo "Checking MysticFunds system status..."
	@echo "===================================="
	@echo ""
	@echo "Service Status:"
	@echo "- Auth Service:   $$(pgrep -f 'auth-service' > /dev/null && echo 'Running' || echo 'Stopped')"
	@echo "- Wizard Service: $$(pgrep -f 'wizard-service' > /dev/null && echo 'Running' || echo 'Stopped')"
	@echo "- Mana Service:   $$(pgrep -f 'mana-service' > /dev/null && echo 'Running' || echo 'Stopped')"
	@echo "- API Gateway:    $$(pgrep -f 'api-gateway' > /dev/null && echo 'Running' || echo 'Stopped')"
	@echo ""
	@echo "Database Status:"
	@for service in $(SERVICES); do \
		result=$$($(SET_ENV) $(PSQL) -h $(PG_HOST) -p $(PG_PORT) -U $(PG_USER) -d $$service -c "SELECT 1" 2>/dev/null && echo "Connected" || echo "Not connected"); \
		printf "- %-12s: %s\n" "$$service" "$$result"; \
	done

logs:
	@echo "Viewing recent logs..."
	@echo "====================="
	@echo ""
	@for service in auth wizard mana; do \
		if [ -f "logs/$$service-service.log" ]; then \
			echo "--- $$service Service (last 10 lines) ---"; \
			tail -n 10 "logs/$$service-service.log"; \
			echo ""; \
		fi; \
	done
	@if [ -f "logs/api-gateway.log" ]; then \
		echo "--- API Gateway (last 10 lines) ---"; \
		tail -n 10 "logs/api-gateway.log"; \
	fi

dev: build
	@echo "Development mode - Auto-restart services..."
	@echo "========================================="
	@echo ""
	@echo "  █▀▄ █▀▀ █ █ █▄█ █▀█ █▀▄ █▀▀"
	@echo "  █▄▀ █▄▄ ▀▄▀ █▀█ █▄█ █▄▀ █▄▄"
	@echo ""
	@echo "========================================="
	@echo "Use 'make stop' to stop all services"
	@$(MAKE) -s start

help:
	@echo "MysticFunds - Available Commands:"
	@echo "================================="
	@echo ""
	@echo "Quick Start:"
	@echo "  make start            - Start entire system"
	@echo "  make stop             - Stop all services"
	@echo "  make status           - Check system status"
	@echo "  make dev              - Development mode"
	@echo "  make logs             - View recent logs"
	@echo ""
	@echo "Building:"
	@echo "  make build            - Build all services and API gateway"
	@echo "  make build-{service}  - Build specific service (e.g., make build-mana)"
	@echo "  make build-gateway    - Build API gateway"
	@echo ""
	@echo "Running (Manual):"
	@echo "  make run             - Run all services (requires multiple terminals)"
	@echo "  make run-{service}   - Run specific service (e.g., make run-mana)"
	@echo "  make run-gateway     - Run API gateway"
	@echo ""
	@echo "Database:"
	@echo "  make create-dbs      - Create databases for all services"
	@echo "  make migrate-up      - Run all migrations up"
	@echo "  make migrate-down    - Run all migrations down"
	@echo "  make migration-status - Check migration status"
	@echo "  make nuke            - Drop all databases (destructive!)"
	@echo ""
	@echo "Development:"
	@echo "  make proto           - Generate protobuf code"
	@echo "  make test            - Run all tests"
	@echo "  make clean           - Clean all service binaries and gateway"
	@echo "  make all             - Create DBs, run migrations, and build (default)"
	@echo ""
	@echo "Deployment:"
	@echo "  make build-railway   - Build all services for Railway deployment"
	@echo "  make start-production - Start services in production mode"
	@echo ""
	@echo "Tip: Use 'make start' for one-command startup!"

# Production build for Railway
build-railway:
	@echo "Building all services for Railway..."
	@for service in $(SERVICES); do \
		echo "Building $$service..."; \
		cd cmd/$$service-service && $(MAKE) build && cd ../..; \
	done
	@echo "Building API Gateway..."
	@cd cmd/$(GATEWAY) && $(MAKE) build && cd ../..
	@echo "Railway build complete!"

# Alias for build-railway
build-all: build-railway

# Start services for production (Railway)
start-production:
	@echo "Starting MysticFunds in production mode..."
	@./railway-start.sh