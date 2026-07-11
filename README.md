# Nail Booking API

Go REST API learning project using Gin, GORM, and PostgreSQL.

## Run Locally

```bash
go run ./cmd/api
```

## Test

```bash
go test ./...
```

## Docker

```bash
docker compose up --build
```

## Structure

```txt
cmd/api        application entrypoint
internal       app-specific packages
pkg            reusable helper packages
docs           project docs
```
