FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY pkg/ pkg/
COPY services/chat/ services/chat/
WORKDIR /app/services/chat
RUN go mod download
RUN CGO_ENABLED=0 go build -o /chat ./cmd/
FROM alpine:3.19
COPY --from=builder /chat /chat
EXPOSE 8094 9094
CMD ["/chat"]
