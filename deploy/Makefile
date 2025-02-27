.PHONY: docker-build docker-run docker-clean start-mongodb

# Detect operating system for MongoDB startup script
ifeq ($(OS),WinWindows_NT)
    MONGODB_START = mongo_scripts/windows.bat
else
    UNAME_S := $(shell uname -s)
    ifeq ($(UNAME_S),Darwin)
        MONGODB_START = mongo_scripts/macos.sh
    else
        MONGODB_START = mongo_scripts/linux.sh
    endif
endif

# Build API Docker image
docker-build-api:
	docker build -t auth-api -f deploy/Dockerfile.main .

# Build Worker Docker image
docker-build-worker:
	docker build -t auth-worker -f deploy/Dockerfile.worker .

# Build both
docker-build: docker-build-api docker-build-worker

# Run API container (foreground)
docker-run-api:
	docker run -p 8080:8080 --name auth-api auth-api

# Run Worker container (detached)
docker-run-worker:
	docker run -d --name auth-worker auth-worker

# Run both (API in foreground, Worker in background)
docker-run:
	$(MAKE) -f deploy/Makefile docker-stop
	$(MAKE) -f deploy/Makefile docker-run-worker
	$(MAKE) -f deploy/Makefile docker-run-api

# Stop and remove running containers
docker-stop:
	docker stop auth-api auth-worker || true
	docker rm auth-api auth-worker || true

# Remove old Docker images
docker-clean: docker-stop
	docker rmi auth-api auth-worker || true

# Start MongoDB locally (if needed)
start-mongodb:
	$(MONGODB_START)
