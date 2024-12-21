FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

WORKDIR /app/cmd

RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o main .

# Final stage
FROM scratch

# Copy SSL certificates for HTTPS (needed for Discord API)
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the binary
COPY --from=builder /app/cmd/main /main

# Run the binary
ENTRYPOINT ["/main"]
