# Detect the operating system
ifeq ($(OS),Windows_NT)
    BINARY_EXT := .exe
    RM := del /Q
    MKDIR := mkdir
else
    BINARY_EXT :=
    RM := rm -f
    MKDIR := mkdir -p
endif

SERVICE_NAME := api-gateway
BINARY_NAME := $(SERVICE_NAME)$(BINARY_EXT)
BINARY_DIR := bin

.PHONY: build run clean test

build:
	@echo "Building $(SERVICE_NAME)..."
	@$(MKDIR) $(BINARY_DIR) 2>/dev/null || true
	@go build -o $(BINARY_DIR)/$(BINARY_NAME) .

run: build
	@echo "Running $(SERVICE_NAME)..."
	@cd $(BINARY_DIR) && ./$(BINARY_NAME)

clean:
	@echo "Cleaning $(SERVICE_NAME)..."
	@$(RM) $(BINARY_DIR)/$(BINARY_NAME) 2>/dev/null || true

test:
	@echo "Testing $(SERVICE_NAME)..."
	@go test ./...