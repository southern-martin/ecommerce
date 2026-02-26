# syntax=docker/dockerfile:1.7

FROM golang:1.26-alpine AS builder

WORKDIR /src

# SERVICE_CMD makes this Dockerfile reusable for each service command
# Example values: cmd/catalog-api, cmd/identity-api, cmd/vendor-api
ARG SERVICE_CMD=cmd/catalog-api

COPY go.mod ./
COPY third_party ./third_party
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /out/service ./${SERVICE_CMD}

FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app
COPY --from=builder /out/service /app/service
COPY --from=builder /src/db/migrations /app/db/migrations

ENV PORT=8080
ENV MIGRATIONS_DIR=/app/db/migrations
EXPOSE 8080

ENTRYPOINT ["/app/service"]
