# Use lightweight Golang image
FROM golang:1.19

WORKDIR /auth-api

# Copy go.mod and go.sum first (for better caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project
COPY . ./

# Generate Swagger docs AFTER copying the source code
RUN go install github.com/swaggo/swag/cmd/swag@latest && \
    swag init -g cmd/api/main.go

# Build the Go binary
RUN go build -o main ./cmd/api/main.go

# Expose the required port
EXPOSE 8080

# Run the compiled binary
CMD ["/app/main"]