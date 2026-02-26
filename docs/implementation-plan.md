# Implementation Plan

## Current status

- `catalog-service` clean architecture scaffold: done
- category/attribute/variant APIs (in-memory): done
- query endpoints for category tree/products: done
- Dockerization baseline per service: done
- PostgreSQL repositories + runtime switch + auto migration: done
- Kong gateway integration (DB-less profile): done

## Next phases

### 1) Dockerization for each service

Objective: every microservice can be containerized with a consistent pattern.

Deliverables:

- Reusable root `Dockerfile` with `SERVICE_CMD` build arg.
- Standard `.dockerignore`.
- Runtime stack `docker-compose.yml` (service + dependencies).
- Microservice compose template for future services in `deploy/`.

Definition of done:

- Service can be built with:
  - `docker build --build-arg SERVICE_CMD=cmd/catalog-api -t ecommerce/catalog-api:dev .`
- Service stack can be started with:
  - `docker compose up --build`

### 2) PostgreSQL repositories

Objective: replace in-memory adapters in production runtime.

Deliverables:

- `internal/infra/postgres` repositories implementing current ports.
- database connection config/env handling.
- migration runner integration and startup documentation.

### 3) API gateway integration

Objective: route customer/vendor/admin traffic through Kong.

Deliverables:

- Kong declarative config (`kong.yml`) with routes and base plugins.
- docker compose profile including Kong for local development.

### 4) Service split execution (current phase)

Objective: expand from catalog into bounded microservices.

Deliverables:

- new `cmd/<service>-api` entrypoints with shared Docker pattern.
- dedicated repo implementations and service-specific migrations.
