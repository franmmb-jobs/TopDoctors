# Build stage
FROM golang:1.24-alpine AS builder

# Install build dependencies for CGO (required by sqlite)
RUN apk add --no-cache build-base

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -o api ./cmd/api/main.go

# Run stage
FROM alpine:latest

RUN apk add --no-cache ca-certificates libc6-compat

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/api .
# Copy configurations
COPY --from=builder /app/configs ./configs

# Expose the API port
EXPOSE 8050

# Run the application
CMD ["./api", "-config=configs/config.docker.yml"]
