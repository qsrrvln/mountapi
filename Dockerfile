# Build stage
FROM golang@sha256:ac09a5f469f307e5da71e766b0bd59c9c49ea460a528cc3e6686513d64a6f1fb AS builder

# Install git, ca-certificates, and tzdata. 
# apk update is removed to avoid extra caching.
RUN apk add --no-cache git ca-certificates tzdata

# Create a non-root user
RUN adduser -D -g '' appuser

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies.
RUN go mod download

# Copy the source code (filtered by .dockerignore)
COPY . .

# Build the Go app statically
# -ldflags="-w -s": Strip debug information for smaller binary
# CGO_ENABLED=0: Disable CGO for static binary
# GOOS=linux: Ensure linux target
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /go/bin/main ./cmd/api

# Run stage (Scratch)
FROM scratch

# Copy timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy CA certificates
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the non-root user
COPY --from=builder /etc/passwd /etc/passwd

# Copy the static binary
COPY --from=builder /go/bin/main /main

# Use the non-root user
USER appuser

# Expose port 8080
EXPOSE 8080

# Command to run the executable
ENTRYPOINT ["/main"]
