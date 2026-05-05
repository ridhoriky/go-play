# ─── Configuration ───────────────────────────────────────────────────────────

APP_NAME = go-play
MAIN_PATH = ./src/cmd/server/main.go
BUILD_DIR = ./bin
CONFIG_FILE = config.yaml

# DB Configuration (adjust as needed or use environment variables)
DB_USER ?= postgres
DB_PASS ?= postgres
DB_HOST ?= localhost
DB_PORT ?= 5432
DB_NAME ?= kasir_db
DB_URL = postgres://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

# ─── Targets ─────────────────────────────────────────────────────────────────

.PHONY: all build run test clean tidy swag migrate-up migrate-down help

all: build

## build: Build the binary
build:
	@echo "Building $(APP_NAME)..."
	@go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)

## run: Run the application
run:
	@go run $(MAIN_PATH)

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

## install-tools: Install necessary tools (golangci-lint, swag, migrate)
install-tools:
	@echo "Installing golangci-lint..."
	@curl curl -sSfL "https://golangci-lint.run/install.sh"
	@sudo mv ./bin/golangci-lint /usr/local/bin/
	@echo "Installing swag..."
	@go install github.com/swaggo/swag/cmd/swag@latest
	@echo "Installing migrate..."
	@go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@echo "All tools installed."	
	@echo "Installation completed."
