# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build mana service
WORKDIR /app/cmd/mana-service
RUN go build -o mana .

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary and config
COPY --from=builder /app/cmd/mana-service/mana .
COPY --from=builder /app/cmd/mana-service/config.yaml ./config.yaml

EXPOSE 50053

CMD ["./mana"]