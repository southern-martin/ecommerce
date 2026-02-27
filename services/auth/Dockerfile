FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy shared pkg module
COPY pkg/ ./pkg/

# Copy auth service
COPY services/auth/ ./services/auth/

WORKDIR /app/services/auth

# Download dependencies
RUN go mod download

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /auth-service ./cmd/main.go

# Final stage
FROM gcr.io/distroless/static-debian12

COPY --from=builder /auth-service /auth-service

EXPOSE 8090 9090

ENTRYPOINT ["/auth-service"]
