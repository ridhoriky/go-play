# KasirApp - Golang REST API

A high-performance, structured RESTful API built with Go for a cashier management system.

---

## 🚀 Features

- **Clean Architecture**: Organized into Handlers, Services, and Repositories.
- **Context-aware Logging**: Structured JSON/Console logging with `zerolog`, featuring unique `request_id` and `component` tracing.
- **Centralized Error Handling**: Uniform API responses and automated error logging via custom middleware.
- **JWT Authentication**: Secure access and refresh token management.
- **Rate Limiting**: Integrated protection against brute-force and spam.
- **Swagger Documentation**: Automatically generated API documentation.
- **Graceful Shutdown**: Handles OS signals for safe database and server closure.

---

## 🛠 Tech Stack

- **Framework**: [Gin Gonic](https://gin-gonic.com/)
- **Database**: [PostgreSQL](https://www.postgresql.org/) with [sqlx](https://github.com/jmoiron/sqlx)
- **Logging**: [Zerolog](https://github.com/rs/zerolog)
- **Migrations**: [golang-migrate](https://github.com/golang-migrate/migrate)
- **Config**: [cleanenv](https://github.com/ilyakaznacheev/cleanenv)

---

## 📋 Prerequisites

- **Go**: v1.25.6 or later
- **PostgreSQL**: v17 or later
- **golang-migrate**: For database migrations
- **swag**: For generating documentation (`go install github.com/swaggo/swag/cmd/swag@latest`)

---

## ⚙️ Setup & Configuration

1. **Clone the repository**
2. **Configure `config.yaml`**:
   Adjust the database and server settings to match your local environment.

3. **Install Dependencies**:

   ```bash
   make tidy
   ```

4. **Run Migrations**:

   ```bash
   make migrate-up
   ```

---

## 🎮 Running the Application

You can use the provided `Makefile` for common tasks:

- **Start Server**: `make run`
- **Build Binary**: `make build`
- **Generate Swagger**: `make swag`
- **Run Tests**: `make test`

The server will be available at `http://localhost:8080`.

---

## 📂 Project Structure

- `src/cmd/server`: Application entry point.
- `src/internal/app`: Application bootstrapping and DI.
- `src/internal/handlers/rest`: HTTP handlers and routing.
- `src/internal/handlers/rest/middleware`: Custom Gin middlewares.
- `src/internal/services`: Business logic layer.
- `src/internal/repositories`: Data access layer (SQL queries).
- `src/internal/models`: DTOs and Entities.
- `etc/migrations`: SQL migration files.

---

## 📝 API Documentation

Once the server is running, you can access the Swagger UI at:
`http://localhost:8080/swagger/index.html`

---
