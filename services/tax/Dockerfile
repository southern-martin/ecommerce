FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy pkg dependency
COPY pkg/ /app/pkg/

# Copy service files
COPY services/tax/go.mod services/tax/go.sum ./services/tax/
WORKDIR /app/services/tax
RUN go mod download

COPY services/tax/ .
RUN CGO_ENABLED=0 GOOS=linux go build -o /tax-service ./cmd/main.go

FROM alpine:3.20

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app
COPY --from=builder /tax-service .

EXPOSE 8098 9098

CMD ["./tax-service"]
