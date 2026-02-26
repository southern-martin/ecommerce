FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY pkg/ pkg/
COPY services/affiliate/ services/affiliate/
WORKDIR /app/services/affiliate
RUN go mod download
RUN CGO_ENABLED=0 go build -o /affiliate ./cmd/
FROM alpine:3.19
COPY --from=builder /affiliate /affiliate
EXPOSE 8097 9097
CMD ["/affiliate"]
