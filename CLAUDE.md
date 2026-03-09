# Ecommerce Platform

Full-stack ecommerce: 19 Go microservices + React/Flutter frontends + Kong Gateway.

## Build & Run

```bash
# Build
make build                       # all 19 services to bin/
make build SVC=product           # single service

# Test
make test                        # all (with race detection)
make test SVC=product            # single service
GOWORK=off go test ./tests/...   # integration/E2E tests (not in go.work)

# Lint & Format
make verify-all                  # fmt + build + vet + lint + test
make vet                         # go vet
make lint                        # golangci-lint
make fmt                         # gofmt

# Infrastructure
make run-infra                   # postgres, redis, nats, elasticsearch, minio
make run-all                     # full stack (docker-compose)
make migrate-up                  # run migrations
make integration-test            # integration tests
make e2e-test                    # E2E tests
```

### Frontend

```bash
cd apps/web && npm install && npm run dev    # React SPA (port 28080)
```

## Go Workspace

- `go.work` at root, Go 1.25.0, 19 services + `pkg/`
- GORM ORM (with pgx/v5 driver) — not raw SQL
- `tests/` module NOT in go.work — use `GOWORK=off`

## Architecture (per service)

```
services/<name>/
├── cmd/main.go
└── internal/
    ├── domain/              # entities, value objects, repository interfaces
    ├── usecase/             # business logic
    ├── adapter/
    │   ├── http/            # Gin handlers, router, DTOs
    │   ├── grpc/            # gRPC server + proto generated code
    │   └── postgres/        # GORM repository implementations
    └── infrastructure/
        ├── config/          # env-based config
        ├── database/        # DB connection, migrations
        └── nats/            # publisher.go, subscriber.go
```

## Code Conventions

### HTTP Layer
- **Gin** framework (not net/http.ServeMux)
- Kong extracts JWT → sets `X-User-ID` header → `ExtractUserID()` middleware reads it
- Middleware chain: CorrelationID → CORS → RateLimit → ExtractUserID → ErrorHandler → routes

### Database
- GORM ORM with PostgreSQL 16, database-per-service (19 schemas)
- Migrations: embed.FS SQL in `services/<name>/migrations/`

### Error Handling
- 6 standardized types via `pkg/errors`: NotFound, Validation, Unauthorized, Forbidden, Conflict, Internal
- `ErrorHandler()` middleware maps these to HTTP status codes

### Observability
- zerolog (structured JSON logging) via `pkg/logger`
- OpenTelemetry + Jaeger (distributed tracing) via `pkg/tracing`
- Prometheus metrics via `pkg/metrics`
- Correlation ID propagated via `X-Request-ID` header

### Events
- NATS JetStream via `pkg/events` (104 subject constants)
- 8/19 services have subscribers wired; 11 remaining

### Auth
- Kong JWT plugin validates tokens globally
- `pkg/auth` for token utilities
- `pkg/middleware.RequireAuth()` ensures auth header present

## Shared Packages (`pkg/`)

auth, cache, circuitbreaker, currency, errors, events, i18n, logger, metrics,
middleware, money, pagination, server, tax, tracing, unitofwork, validator

## Docker Ports (28xxx/29xxx range)

| Service | HTTP | gRPC | | Service | HTTP | gRPC |
|---------|------|------|-|---------|------|------|
| kong | 28000 | — | | ai | 28092 | 29092 |
| auth | 28090 | 29090 | | promotion | 28093 | 29093 |
| user | 28091 | 29091 | | return | 28094 | 29094 |
| product | 28081 | 29081 | | shipping | 28095 | 29095 |
| cart | 28082 | 29082 | | loyalty | 28096 | 29096 |
| order | 28083 | 29083 | | affiliate | 28097 | 29097 |
| payment | 28084 | 29084 | | tax | 28098 | 29098 |
| search | 28085 | 29085 | | cms | 28099 | 29099 |
| review | 28086 | 29086 | | web | 28080 | — |
| notification | 28087 | 29087 | | postgres | 15432 | — |
| chat | 28088 | 29088 | | redis | 16379 | — |
| media | 28089 | 29089 | | nats | 14222 | — |

## Git Workflow

- Feature branch → merge to `develop` with `--no-ff`
- Never push directly to main (releases only)
- `gh` CLI is NOT authenticated on this machine

## Frontend

- React 18.3 + Vite + TypeScript + Tailwind + Shadcn/ui
- TanStack Query v5 (server state), Zustand (client state)
- Flutter mobile apps: `apps/mobile/buyer/`, `apps/mobile/seller/`

## Key Files

- Architecture: `ARCHITECTURE_PLAN.md` (4,475 lines)
- Makefile: 246 lines, all build/test/deploy targets
- Docker: `docker-compose.yml` (full stack), `docker-compose.override.yml` (observability)
- Kong: `kong/kong.yml` (DB-less declarative config)
- K8s: `k8s/` (HPA, PDB, monitoring)
- CI/CD: `.github/workflows/` (go-ci, react-ci, flutter-ci)
