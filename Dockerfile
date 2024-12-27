FROM golang:1.23-alpine

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

WORKDIR /app/cmd

RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Run the binary
ENTRYPOINT ["/app/cmd/main"]
