# Nail Booking API

Go REST API learning project using Gin, GORM, and PostgreSQL.

## Run Locally

```bash
go run ./cmd/api
```

## Keep Supabase Active

Use this endpoint from a cron job or uptime monitor to create a lightweight database activity:

```bash
curl https://your-api-domain.com/api/keep-alive
```

The endpoint runs `SELECT 1` against PostgreSQL and returns `200 OK` when the database is reachable.

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
