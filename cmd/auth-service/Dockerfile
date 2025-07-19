# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build auth service
WORKDIR /app/cmd/auth-service
RUN go build -o auth .

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary and config
COPY --from=builder /app/cmd/auth-service/auth .
COPY --from=builder /app/cmd/auth-service/config.yaml ./config.yaml

EXPOSE 50051

CMD ["./auth"]