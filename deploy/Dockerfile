# Use lightweight Golang image
FROM golang:1.19

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o main ./cmd/api/main.go

EXPOSE 8080
CMD ["/app/main"]
