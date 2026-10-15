# =========================
# Build stage
# =========================
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy dependency files first
COPY go.mod go.sum ./

RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o auth-service ./cmd/server


# =========================
# Runtime stage
# =========================
FROM alpine:3.22

WORKDIR /app

# Install certificates for HTTPS/TLS communication
RUN apk --no-cache add ca-certificates

COPY --from=builder /app/auth-service .

EXPOSE 8080

CMD ["./auth-service"]