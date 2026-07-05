# ─── Configuration ───────────────────────────────────────────────────────────

APP_NAME = go-play
MAIN_PATH = ./src/cmd/server/main.go
BUILD_DIR = ./bin
CONFIG_FILE = config.yaml

# Load environment variables from .env file if it exists
ifneq (,$(wildcard .env))
    include .env
    export
endif

# DB Configuration (adjust as needed or use environment variables)
DB_USER ?= postgres
DB_PASS ?= mypassword
DB_HOST ?= localhost
DB_PORT ?= 5432
DB_NAME ?= kasir_db
DB_URL = postgres://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

SMTP_USERNAME ?=
SMTP_PASSWORD ?=
SMTP_SENDER_EMAIL ?=

# ─── Targets ─────────────────────────────────────────────────────────────────

.PHONY: all build run test clean tidy swag migrate-up migrate-down seed help lint install-tools lefthook-run lefthook-uninstall 

all: build

## build: Build the binary
build:
	@echo "Building $(APP_NAME)..."ocker-build docker-up docker-down docker-logs docker-ps docker-stop-all
	@go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)

## run: Run the application
run:
	@SMTP_USERNAME="$(SMTP_USERNAME)" SMTP_PASSWORD="$(SMTP_PASSWORD)" SMTP_SENDER_EMAIL="$(SMTP_SENDER_EMAIL)" go run $(MAIN_PATH)

## test: Run tests
test:
	@go test -v ./...

## clean: Remove build artifacts and logs
clean:
	@rm -rf $(BUILD_DIR)
	@rm -rf ./logs/*.log
	@echo "Cleanup completed."

## tidy: Run go mod tidy
tidy:
	@go mod tidy

## swag: Generate Swagger documentation
swag:
	@swag init -g $(MAIN_PATH) -o ./docs --parseDependency --parseInternal

## migrate-up: Run database migrations up
migrate-up:
	@migrate -path ./etc/migrations -database "$(DB_URL)" up

## migrate-down: Run database migrations down
migrate-down:
	@migrate -path ./etc/migrations -database "$(DB_URL)" down 1

## seed: Run database seeder
seed:
	@go run ./src/cmd/seeder/main.go

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^##//p' Makefile | column -t -s ':' |  sed -e 's/^/ /'

## lint: Run golangci-lint
lint:
	@echo "Running golangci-lint..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
	    golangci-lint run; \
	else \
	    echo "golangci-lint not installed. Run: make install-tools"; \
	    exit 1; \
	fi
	@echo "golangci-lint finished."

## install-tools: Install all tools (linter, lefthook, swag, migrate, hey) and setup hooks
install-tools:
	@echo "Installing golangci-lint..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Installing lefthook..."
	@go install github.com/evilmartians/lefthook/v2@latest
	@echo "Installing swag..."
	@go install github.com/swaggo/swag/cmd/swag@latest
	@echo "Installing migrate..."
	@go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@echo "Installing Lefthook Git hooks..."
	@lefthook install
	@echo "All tools installed and hooks activated."

## lefthook-run: Run Lefthook pre-commit manually (for testing)
lefthook-run:
	@lefthook run pre-commit

## lefthook-uninstall: Remove Lefthook hooks
lefthook-uninstall:
	@lefthook uninstall

## generate-keys: Generate ECDSA P-256 key pairs for JWT (ES256)
generate-keys:
	@chmod +x scripts/generate-jwt-keys.sh
	@./scripts/generate-jwt-keys.sh

