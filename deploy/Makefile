.PHONY: build run docker-build start-mongodb

# Detect operating system
ifeq ($(OS),Windows_NT)
    MONGODB_START = mongo_scripts/windows.bat
else
    UNAME_S := $(shell uname -s)
    ifeq ($(UNAME_S),Darwin)
        MONGODB_START = mongo_scripts/macos.sh
    else
        MONGODB_START = mongo_scripts/linux.sh
    endif
endif

build:
	go mod tidy
	swag init -g cmd/api/main.go
	go build -o bin/api ./cmd/api/main.go

run: build
	./bin/api

docker-build:
	docker build -t auth-api .

start-mongodb:
	$(MONGODB_START)