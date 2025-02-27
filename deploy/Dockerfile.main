# Use official Go image
FROM golang:1.22.4-alpine

# Set the working directory inside the container
WORKDIR /app

# Install dependencies (Make, Swag, MongoDB tools)
RUN apk add --no-cache make git

# Copy everything into the container
COPY . .

# Install Go dependencies
RUN go mod tidy

# Generate Swagger docs
RUN go install github.com/swaggo/swag/cmd/swag@latest && swag init -g cmd/api/main.go

# Build the API
RUN go build -o bin/api ./cmd/api/main.go

# Expose the API port
EXPOSE 8080

# Start the API when the container runs
CMD ["./bin/api"]
