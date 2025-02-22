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

# Build Docker image
docker-build:
	docker build -t auth-api -f deploy/Dockerfile .

# Run Docker container
docker-run:
	docker run -p 8080:8080 auth-api

# Remove old Docker image (for rebuilding)
docker-clean:
	docker rmi auth-api

# Start MongoDB locally (if needed)
start-mongodb:
	$(MONGODB_START)
