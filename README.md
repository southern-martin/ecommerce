# Ecommerce Marketplace Platform

A full-stack B2C marketplace platform built with **Go microservices** (Clean Architecture) and a **React SPA** frontend. Connects buyers and sellers with AI-powered features including semantic search, image search, and AI-generated product descriptions.

## Architecture

```
Client (Browser/Mobile)
        │
   Kong API Gateway (:8000)
        │
   ┌────┴─────────────────────────────────┐
   │    19 Go Microservices (REST + gRPC)  │
   │    Inter-service: gRPC + NATS events  │
   └────┬─────────────────────────────────┘
        │
   PostgreSQL │ Redis │ Elasticsearch │ MinIO
```

Each service follows **Clean Architecture** (domain → use case → adapter → infrastructure) with its own PostgreSQL database. Services communicate via gRPC for synchronous calls and NATS JetStream for async events.

## Services

| Service | HTTP | gRPC | Description |
|---------|------|------|-------------|
| auth | 8090 | 9090 | Registration, login, JWT, OAuth |
| user | 8091 | 9091 | Profiles, addresses, seller verification |
| product | 8081 | 9081 | Catalog, categories, variants, attributes |
| cart | 8082 | 9082 | Shopping cart (Redis-backed) |
| order | 8083 | 9083 | Order lifecycle, status tracking |
| payment | 8084 | 9084 | Stripe integration, refunds |
| search | 8085 | 9085 | Elasticsearch-powered product search |
| review | 8086 | 9086 | Ratings, reviews, moderation |
| notification | 8087 | 9087 | Email, push, WebSocket notifications |
| chat | 8088 | 9088 | Buyer-seller messaging |
| media | 8089 | 9089 | Image/file upload to S3/MinIO |
| ai | 8092 | 9092 | AI descriptions, semantic search, recommendations |
| promotion | 8093 | 9093 | Coupons, vouchers, flash sales |
| return | 8094 | 9094 | Returns, refunds, disputes |
| shipping | 8095 | 9095 | Carrier integration, tracking |
| loyalty | 8096 | 9096 | Points, cashback, membership tiers |
| affiliate | 8097 | 9097 | Referral tracking, commissions |
| tax | 8098 | 9098 | Tax rules engine, jurisdiction config |
| cms | 8099 | 9099 | Banners, landing pages, content scheduling |

## Tech Stack

| Layer | Technology |
|-------|------------|
| Backend | Go (Gin), Clean Architecture, GORM |
| Frontend | React 18 + TypeScript + Vite + Tailwind + Shadcn/ui |
| State | Zustand (client) + TanStack Query (server) |
| API Gateway | Kong (DB-less, declarative YAML) |
| Database | PostgreSQL 16 + pgvector |
| Cache | Redis 7 |
| Search | Elasticsearch 8 |
| Messaging | NATS JetStream |
| Storage | S3-compatible (MinIO for dev) |
| Auth | JWT + OAuth2 (Google, Facebook, Apple) |
| Payments | Stripe + Stripe Connect |
| Mobile | Flutter 3.x (Dart) |
| Observability | OpenTelemetry + Jaeger, Prometheus + Grafana |

## Project Structure

```
ecommerce/
├── apps/
│   ├── web/                  # React SPA (Vite + TypeScript)
│   └── mobile/               # Flutter apps (Buyer + Seller/Admin)
├── services/
│   ├── auth/                 # Each service follows:
│   ├── product/              #   cmd/          - entrypoint
│   ├── order/                #   internal/
│   ├── ...                   #     domain/     - entities, value objects
│   └── ai/                   #     usecase/    - application logic
│                             #     port/       - interfaces
│                             #     adapter/    - HTTP/gRPC handlers
│                             #     infrastructure/ - DB, cache, messaging
├── pkg/                      # Shared Go packages
│   ├── auth/                 #   JWT utilities
│   ├── errors/               #   Standardized error handling
│   ├── events/               #   NATS event publishing
│   ├── logger/               #   Structured logging (zerolog)
│   ├── middleware/            #   Common HTTP middleware
│   ├── server/               #   HTTP + gRPC server bootstrap
│   ├── validator/            #   Request validation
│   └── ...                   #   currency, i18n, money, tax, etc.
├── proto/                    # Protobuf definitions (per service)
├── kong/                     # Kong gateway config (kong.yml)
├── infra/                    # Docker, Terraform, K8s manifests
├── scripts/                  # Dev utilities
├── docker-compose.yml        # Full local stack
├── go.work                   # Go workspace (all services + pkg)
├── Makefile                  # Build, test, lint, deploy targets
└── ARCHITECTURE_PLAN.md      # Full architecture documentation
```

## Quick Start

### Prerequisites

- Go 1.25+
- Docker & Docker Compose
- Node.js 18+ & npm

### 1. Start Infrastructure

```bash
# Start PostgreSQL, Redis, NATS, Elasticsearch, MinIO
make run-infra
```

### 2. Run Backend Services

```bash
# Build all services
make build

# Or run individual services
cd services/auth && go run ./cmd/
```

### 3. Run Frontend

```bash
cd apps/web
npm install
npm run dev
# Open http://localhost:3000
```

### 4. Full Docker Stack

```bash
docker compose up -d
```

## Infrastructure Ports (Docker)

| Service | Host Port | Container Port |
|---------|-----------|----------------|
| PostgreSQL | 15432 | 5432 |
| Redis | 16379 | 6379 |
| NATS | 14222 | 4222 |
| Elasticsearch | 19200 | 9200 |
| MinIO | 19000 | 9000 |
| Kong Proxy | 8000 | 8000 |
| Kong Admin | 8001 | 8001 |

## Make Targets

```
make build          # Build all Go services
make test           # Run all tests with race detection
make lint           # Run golangci-lint
make proto          # Generate protobuf code
make run-infra      # Start infrastructure containers
make run-all        # Start everything via docker-compose
make migrate-up     # Run database migrations
make docker-build   # Build Docker images
make verify-all     # Format + build + vet + lint + test
```

## Environment Variables

Copy `.env.example` to `.env` and adjust values:

```bash
cp .env.example .env
```

Key variables: `POSTGRES_*`, `REDIS_URL`, `NATS_URL`, `JWT_SECRET`, `S3_*`, `ELASTICSEARCH_URL`

## License

Private repository.
