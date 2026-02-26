FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY pkg/ ./pkg/
COPY services/product/ ./services/product/

WORKDIR /app/services/product

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/product-service ./cmd/main.go

FROM alpine:3.20

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/product-service .

EXPOSE 8081 9081

CMD ["./product-service"]
