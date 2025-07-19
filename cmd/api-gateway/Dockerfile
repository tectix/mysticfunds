# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build API gateway
WORKDIR /app/cmd/api-gateway
RUN go build -o api-gateway .

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates curl

WORKDIR /root/

# Copy binary, config, and web assets
COPY --from=builder /app/cmd/api-gateway/api-gateway .
COPY --from=builder /app/cmd/api-gateway/config.yaml ./config.yaml
COPY --from=builder /app/web ./web/

EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8080/health || exit 1

CMD ["./api-gateway"]