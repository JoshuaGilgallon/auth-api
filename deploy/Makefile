.PHONY: build run docker-build

# Detect operating system
ifeq ($(OS),Windows_NT)
    MONGODB_START = sc start MongoDB
	KUBECT_INSTALL = choco install kubernetes-cli -y
else
    UNAME_S := $(shell uname -s)
    ifeq ($(UNAME_S),Darwin)
        MONGODB_START = brew services start mongodb-community
		KUBECT_INSTALL = brew install kubectl
    else
        MONGODB_START = sudo systemctl start mongod
		KUBECT_INSTALL = idk :sob:
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
