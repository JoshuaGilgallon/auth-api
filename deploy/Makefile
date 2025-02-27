.PHONY: docker-build docker-run docker-clean start-mongodb

# Detect operating system for MongoDB startup script
ifeq ($(OS),Windows_NT)
    MONGODB_START = mongo_scripts/windows.bat
else
    MONGODB_START = $(shell uname -s | grep -q Darwin && echo mongo_scripts/macos.sh || echo mongo_scripts/linux.sh)
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
	-docker stop auth-api auth-worker
	-docker rm auth-api auth-worker

# Remove old Docker images
docker-clean: docker-stop
	-docker rmi auth-api auth-worker

# Start MongoDB locally (if needed)
start-mongodb:
ifeq ($(OS),Windows_NT)
	powershell -ExecutionPolicy Bypass -File $(MONGODB_START)
else
	chmod +x $(MONGODB_START)
	./$(MONGODB_START)
endif
