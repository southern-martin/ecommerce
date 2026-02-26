# Kong Gateway (DB-less)

This configuration runs Kong in declarative mode and proxies catalog endpoints.

## Route mapping

- Incoming: `/api/catalog/*`
- Upstream: `catalog-api:8080/*`

Example:

- `GET /api/catalog/healthz` -> `GET /healthz`
- `POST /api/catalog/categories` -> `POST /categories`

## Enabled baseline plugins

- `cors`
- `rate-limiting`
- `correlation-id`
- `request-size-limiting`
- `prometheus`

## Local endpoints

- Proxy: `http://localhost:8000`
- Admin API (dev only): `http://localhost:8001`
- Metrics: `http://localhost:8100/metrics`
