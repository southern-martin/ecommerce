FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Copy pkg dependency first
COPY pkg/ /app/pkg/

# Copy payment service
COPY services/payment/ /app/services/payment/

WORKDIR /app/services/payment

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /payment-service ./cmd/main.go

FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /payment-service .

EXPOSE 8084 9084

CMD ["./payment-service"]
