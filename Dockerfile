# Build stage
FROM golang:1.23.3-alpine3.19 AS builder

# Set the working directory
WORKDIR /app

# Copy the Go modules manifests
COPY go.mod go.sum ./

# Download the Go modules
RUN go mod download

# Copy the source code
COPY . .

# Build the application binary
RUN go build -o main .

# Runtime stage
FROM gcr.io/distroless/base-debian10

# Set working directory
WORKDIR /

# Copy the application binary
COPY --from=builder /app/main .

# Expose application port
EXPOSE 8080

# Run the application
CMD ["./main"]
