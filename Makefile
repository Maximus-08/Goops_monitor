.PHONY: build run test clean docker docker-run

# Variables
BINARY_NAME=monitor_bin
DOCKER_IMAGE=goops-monitor

# Build the binary
build:
	go build -o $(BINARY_NAME) monitor/*.go

# Run locally
run: build
	./$(BINARY_NAME)

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -f $(BINARY_NAME)
	rm -f monitor.log

# Build Docker image
docker:
	docker build -t $(DOCKER_IMAGE) .

# Run Docker container
docker-run: docker
	docker run -p 8080:8080 -p 8081:8081 $(DOCKER_IMAGE)

# Run with docker-compose
compose-up:
	docker-compose up -d

compose-down:
	docker-compose down

# Show logs
logs:
	docker-compose logs -f
