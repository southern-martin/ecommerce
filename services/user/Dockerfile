FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy shared pkg module
COPY pkg/ ./pkg/

# Copy user service module
COPY services/user/ ./services/user/

WORKDIR /app/services/user

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /user-service ./cmd/main.go

FROM gcr.io/distroless/static-debian12

COPY --from=builder /user-service /user-service

EXPOSE 8091 9091

ENTRYPOINT ["/user-service"]
