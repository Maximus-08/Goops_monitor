# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o monitor_bin monitor/*.go

# Run stage
FROM alpine:latest

LABEL org.opencontainers.image.source="https://github.com/Maximus-08/goops-monitor"
LABEL org.opencontainers.image.description="Lightweight health monitoring tool"
LABEL org.opencontainers.image.version="1.0.0"

RUN apk --no-cache add ca-certificates wget

# Create non-root user for security
RUN adduser -D -g '' appuser

WORKDIR /app

COPY --from=builder /app/monitor_bin .
COPY config.json .

# Change ownership to non-root user
RUN chown -R appuser:appuser /app
USER appuser

# Expose ports for API and demo server
EXPOSE 8081 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8081/status || exit 1

CMD ["./monitor_bin"]
