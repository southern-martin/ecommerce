FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY pkg/ pkg/
COPY services/promotion/ services/promotion/
WORKDIR /app/services/promotion
RUN go mod download
RUN CGO_ENABLED=0 go build -o /promotion ./cmd/
FROM alpine:3.19
COPY --from=builder /promotion /promotion
EXPOSE 8093 9093
CMD ["/promotion"]
