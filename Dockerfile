# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod* ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Runtime stage
FROM scratch

WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Copy the shortcuts file
COPY shortcuts.json .

# Expose port 80
EXPOSE 80

# Expose volume for shortcuts
VOLUME ["/data"]

# Run the application
CMD ["./main", ":80"]