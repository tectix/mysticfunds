# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build wizard service
WORKDIR /app/cmd/wizard-service
RUN go build -o wizard .

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary and config
COPY --from=builder /app/cmd/wizard-service/wizard .
COPY --from=builder /app/cmd/wizard-service/config.yaml ./config.yaml

EXPOSE 50052

CMD ["./wizard"]