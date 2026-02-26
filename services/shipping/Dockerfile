FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY pkg/ ./pkg/
COPY services/shipping/ ./services/shipping/

WORKDIR /app/services/shipping
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/shipping ./cmd/main.go

FROM alpine:3.20

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/shipping .

EXPOSE 8085 9085

CMD ["./shipping"]
