.PHONY: build run docker-build

# Detect operating system
ifeq ($(OS),Windows_NT)
    MONGODB_START = net start MongoDB
else
    UNAME_S := $(shell uname -s)
    ifeq ($(UNAME_S),Darwin)
        MONGODB_START = brew services start mongodb-community
    else
        MONGODB_START = sudo systemctl start mongod
    endif
endif

build:
	go mod tidy
	swag init -g cmd/api/main.go
	go build -o bin/api ./cmd/api/main.go
	$(MONGODB_START)

run: build
	./bin/api

docker-build:
	docker build -t auth-api .
