.PHONY: run build test clean deps docker-build docker-run

# Run the application locally
run:
	go run cmd/api/main.go

# Build the binary
build:
	go build -o bin/api cmd/api/main.go

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/
	go clean

# Download and tidy dependencies
deps:
	go mod download
	go mod tidy

# Build Docker image
docker-build:
	docker build -t event-campus-api:latest .

# Run Docker container
docker-run:
	docker run -p 8080:8080 --env-file .env event-campus-api:latest

# Run with live reload (requires air: go install github.com/cosmtrek/air@latest)
dev:
	air
