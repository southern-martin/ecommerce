FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY pkg/ ./pkg/
COPY services/cart/ ./services/cart/

WORKDIR /app/services/cart

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /cart-service ./cmd/main.go

FROM alpine:3.20

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /cart-service .

EXPOSE 8082 9082

CMD ["./cart-service"]
