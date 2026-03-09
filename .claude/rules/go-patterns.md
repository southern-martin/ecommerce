---
paths:
  - "services/**/*.go"
  - "pkg/**/*.go"
---
# Go Patterns

## Service Structure

```
services/<name>/
├── cmd/main.go                          # entry point, wire deps, start server
└── internal/
    ├── domain/                          # entities, value objects, repo interfaces
    ├── usecase/                         # business logic
    ├── adapter/
    │   ├── http/                        # Gin handlers, router.go, DTOs
    │   ├── grpc/                        # gRPC server + proto code
    │   └── postgres/                    # GORM repository implementations
    └── infrastructure/
        ├── config/config.go             # env-based config loading
        ├── database/                    # DB connection, migrations
        └── nats/                        # publisher.go, subscriber.go
```

## HTTP Layer (Gin)

- Routes defined in `adapter/http/router.go`
- Use `pkg/middleware` for cross-cutting: CorrelationID, CORS, RateLimit, ExtractUserID, ErrorHandler
- Kong sets `X-User-ID` header after JWT validation — read via `ExtractUserID()` middleware

## Database (GORM)

- GORM ORM with pgx/v5 driver — NOT raw SQL
- Database-per-service pattern (19 PostgreSQL schemas)
- Migrations via embed.FS SQL files in `services/<name>/migrations/`

## Error Handling

Use `pkg/errors` standardized types:
```go
errors.NewNotFound("product", id)
errors.NewValidation("price must be positive")
errors.NewConflict("email already exists")
```
`ErrorHandler()` middleware maps these to HTTP status codes automatically.

## NATS Events

- Subject constants in `pkg/events` (104 defined)
- Publisher: `infrastructure/nats/publisher.go`
- Subscriber: `infrastructure/nats/subscriber.go`
- Wire in `cmd/main.go`

## Observability

- `pkg/logger` (zerolog) — use `logger.Info().Str("key", val).Msg("message")`
- `pkg/tracing` — OpenTelemetry spans auto-attached by middleware
- `pkg/metrics` — Prometheus counters/histograms via middleware
