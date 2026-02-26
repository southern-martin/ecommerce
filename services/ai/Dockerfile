FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY pkg/ pkg/
COPY services/ai/ services/ai/
WORKDIR /app/services/ai
RUN go mod download
RUN CGO_ENABLED=0 go build -o /ai ./cmd/

FROM alpine:3.19
COPY --from=builder /ai /ai
EXPOSE 8095 9095
CMD ["/ai"]
