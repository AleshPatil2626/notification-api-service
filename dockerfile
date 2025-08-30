# -----------------------
# Build stage
# -----------------------
FROM golang:1.24-alpine AS builder
WORKDIR /app

# Copy go.mod only (no go.sum yet)
COPY go.mod ./

# Download modules (will skip if no external modules)
RUN go mod download

# Copy the entire source code
COPY . .

# Build the Go binary
RUN go build -o notification-api-service .

# -----------------------
# Final stage
# -----------------------
FROM alpine:latest
WORKDIR /app

# Install bash (and optional tools for convenience)
RUN apk add --no-cache bash curl vim

# Copy everything from builder (source code + binary)
COPY --from=builder /app /app

# Expose the port your Go app listens on
EXPOSE 8006

# Start the Go server
CMD ["./notification-api-service"]
