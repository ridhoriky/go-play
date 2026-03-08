# KasirApp - Golang REST API

KasirApp is a RESTful API built with Go for a simple cashier system using Gin and PostgreSQL.

---

## 🛠 Tech Stack

- Go
- Gin
- PostgreSQL
- sqlx
- golang-migrate

---

## 📋 Requirements

```bash
go version
psql --version
migrate -version
```

Minimum:

- Go go1.25.6
- PostgreSQL 17
- golang-migrate v4.19.1

## ⚙️ Configuration

Edit config.yaml:

```yaml
app:
  name: KasirApp
  port: 8080

database:
  host: localhost
  port: 5432
  user: postgres
  password: postgres
  name: kasir_db
```

### Create database

```sql
createdb kasir_db
```

### 📦 Installation

```bash
go mod tidy
```

### 🗄 Database Migration

Set DB URL:

```bash
export DB_URL="postgres://postgres:postgre@localhost:5432/kasir_db?sslmode=disable"
```

Run migration:

```bash
migrate -path ./etc/migrations -database "$DB_URL" up
```

▶️ Run Application

```bash
go run cmd/server/main.go
```

Server runs at:

```bash
http://localhost:8080
```
