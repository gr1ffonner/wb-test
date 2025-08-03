# Build stage
FROM golang:1.21 AS builder

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go

# Runtime stage
FROM alpine:latest


# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/main .


# Expose port (adjust if your app uses a different port)
EXPOSE 8080

# Run the application
CMD ["./main"] 