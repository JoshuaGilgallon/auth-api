.PHONY: build run docker-build

build:
	go mod tidy
	swag init -g cmd/api/main.go
	go build -o bin/api ./cmd/api/main.go
	brew services start mongodb-community

run: build
	./bin/api

docker-build:
	docker build -t auth-api .
