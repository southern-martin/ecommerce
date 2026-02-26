FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY pkg/ pkg/
COPY services/loyalty/ services/loyalty/
WORKDIR /app/services/loyalty
RUN go mod download
RUN CGO_ENABLED=0 go build -o /loyalty ./cmd/
FROM alpine:3.19
COPY --from=builder /loyalty /loyalty
EXPOSE 8096 9096
CMD ["/loyalty"]
