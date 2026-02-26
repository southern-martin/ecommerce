# Catalog Service

Catalog domain for a multivendor ecommerce platform, implemented with clean architecture boundaries.

## Scope in this scaffold

- Category tree (`parent -> children`)
- Products with one primary category and optional extra categories
- Category-scoped attributes and options
- Product-level attribute values (non-variant specs)
- Variant matrix generation using selected attribute options

## Structure

- `cmd/catalog-api`: HTTP entrypoint
- `db/migrations`: PostgreSQL schema
- `docs/openapi.yaml`: API contract
- `internal/domain`: entities and domain validation
- `internal/port`: repository interfaces
- `internal/usecase`: application use cases
- `internal/infra/memory`: in-memory adapter repositories
- `internal/adapter/httpapi`: HTTP handlers and request/response mapping

## Run tests

```bash
go test ./...
```

## Run API

```bash
go run ./cmd/catalog-api
```

Default address: `http://localhost:8080`

### Runtime mode

- `REPOSITORY_MODE=memory` (default): in-memory repositories
- `REPOSITORY_MODE=postgres`: PostgreSQL repositories

PostgreSQL mode envs:

- `DATABASE_URL` (required)
- `AUTO_MIGRATE` (`true`/`false`, default `true`)
- `MIGRATIONS_DIR` (default `db/migrations`)

## Dockerization

### Build one service image

Use the root `Dockerfile` for every service by changing `SERVICE_CMD`.

```bash
docker build --build-arg SERVICE_CMD=cmd/catalog-api -t ecommerce/catalog-api:dev .
```

For other services, replace `SERVICE_CMD` when their `cmd/<service>-api` entrypoint exists.

### Run local stack

```bash
docker compose up --build
```

This starts:

- `catalog-api` on `:8080`
- `postgres` on `:5432`

The container runs with `REPOSITORY_MODE=postgres` and auto-applies migrations on startup.

### Run with Kong gateway profile

```bash
docker compose --profile gateway up --build
```

Gateway endpoints:

- Proxy base: `http://localhost:8000`
- Admin API (dev): `http://localhost:8001`
- Metrics: `http://localhost:8100/metrics`

Catalog API through Kong:

- `GET http://localhost:8000/api/catalog/healthz`
- `POST http://localhost:8000/api/catalog/categories`
- `GET http://localhost:8000/api/catalog/categories/tree`

### Multi-service template

Use `deploy/docker-compose.services.example.yml` as the baseline to add each microservice with the same Docker pattern.
Use `deploy/kong/kong.yml` as the baseline for gateway routes/plugins.

## Implemented endpoints

- `GET /healthz`
- `POST /categories`
- `GET /categories/tree`
- `GET /categories/{categoryID}/children`
- `GET /categories/{categoryID}/products`
- `POST /categories/{categoryID}/attributes`
- `GET /categories/{categoryID}/attributes`
- `POST /attributes/{attributeID}/options`
- `POST /products`
- `GET /products/{productID}`
- `PUT /products/{productID}/attributes`
- `GET /products/{productID}/attributes`
- `POST /products/{productID}/variants/generate`
- `GET /products/{productID}/variants`
