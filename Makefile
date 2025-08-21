# Variables
BINARY_NAME=analytics-service
DOCKER_IMAGE=analytics-service
DOCKER_TAG=latest

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Build the application
build:
	$(GOBUILD) -o $(BINARY_NAME) ./cmd/server

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME).exe

# Run tests
test:
	$(GOTEST) -v ./...

# Run tests with coverage
test-coverage:
	$(GOTEST) -v -cover ./...

# Download dependencies
deps:
	$(GOMOD) download

# Tidy dependencies
tidy:
	$(GOMOD) tidy

# Run the application
run:
	$(GOBUILD) -o $(BINARY_NAME) ./cmd/server
	./$(BINARY_NAME)

# Run the application (Windows)
run-windows:
	$(GOBUILD) -o $(BINARY_NAME).exe ./cmd/server
	./$(BINARY_NAME).exe

# Build Docker image
docker-build:
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

# Run Docker container
docker-run:
	docker run -p 8080:8080 $(DOCKER_IMAGE):$(DOCKER_TAG)

# Build and run Docker container
docker: docker-build docker-run

# Format code
fmt:
	$(GOCMD) fmt ./...

# Vet code
vet:
	$(GOCMD) vet ./...

# Lint code (requires golangci-lint)
lint:
	golangci-lint run

# Benchmark analytics endpoint
bench:
	@echo "Running analytics benchmark..."
	@curl -X POST http://localhost:8080/analytics \
		-H "Content-Type: application/json" \
		-d @examples/sample.json \
		-w "\nTime: %{time_total}s\n" \
		-s

# Test API (PowerShell)
test-api:
	@echo "Testing API with PowerShell..."
	@powershell -ExecutionPolicy Bypass -File scripts/test_api.ps1

# Load test tokens
load-tokens:
	@echo "Loading tokens from LogPas.txt..."
	@if [ -f routes/LogPas.txt ]; then \
		echo "Tokens loaded successfully"; \
	else \
		echo "Warning: LogPas.txt not found"; \
	fi

# Install development tools (legacy)
install-tools-legacy:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run benchmarks
benchmark:
	@echo "Running benchmarks..."
	@go test -bench=. -benchmem ./internal/analytics/

# Run race detector
test-race:
	@echo "Running tests with race detector..."
	@go test -race -v ./...

# Generate coverage report
coverage:
	@echo "Generating coverage report..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Docker Compose commands
compose-up:
	@echo "Starting services with Docker Compose..."
	@docker-compose up -d

compose-down:
	@echo "Stopping services..."
	@docker-compose down

compose-logs:
	@docker-compose logs -f

# Kubernetes commands
k8s-deploy:
	@echo "Deploying to Kubernetes..."
	@kubectl apply -f k8s/

k8s-delete:
	@echo "Deleting from Kubernetes..."
	@kubectl delete -f k8s/

# Development tools
install-tools:
	@echo "Installing development tools..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/axw/gocov/gocov@latest
	@go install github.com/AlekSi/gocov-xml@latest
	@go install github.com/tebeka/go2xunit@latest

# Security scan
security-scan:
	@echo "Running security scan..."
	@if command -v trivy >/dev/null 2>&1; then \
		trivy fs .; \
	else \
		echo "Trivy not found. Install with: go install github.com/aquasecurity/trivy/cmd/trivy@latest"; \
	fi

# Performance profiling
profile:
	@echo "Starting performance profiling..."
	@go build -o analytics-service ./cmd/server
	@./analytics-service &
	@sleep 2
	@curl -X POST http://localhost:8080/analytics -H "Content-Type: application/json" -d @examples/sample.json
	@curl http://localhost:8080/debug/pprof/profile?seconds=30 -o cpu.prof
	@curl http://localhost:8080/debug/pprof/heap -o memory.prof
	@pkill analytics-service

# Performance profiling (Windows)
profile-windows:
	@echo "Starting performance profiling (Windows)..."
	@go build -o analytics-service.exe ./cmd/server
	@start /B analytics-service.exe
	@timeout /t 2 /nobreak > nul
	@curl -X POST http://localhost:8080/analytics -H "Content-Type: application/json" -d @examples/sample.json
	@curl http://localhost:8080/debug/pprof/profile?seconds=30 -o cpu.prof
	@curl http://localhost:8080/debug/pprof/heap -o memory.prof
	@taskkill /F /IM analytics-service.exe

# Help
help:
	@echo "Available commands:"
	@echo "  build         - Build the application"
	@echo "  clean         - Clean build artifacts"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  test-race     - Run tests with race detector"
	@echo "  coverage      - Generate coverage report"
	@echo "  benchmark     - Run benchmarks"
	@echo "  deps          - Download dependencies"
	@echo "  tidy          - Tidy dependencies"
	@echo "  run           - Build and run the application"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Run Docker container"
	@echo "  docker        - Build and run Docker container"
	@echo "  compose-up    - Start services with Docker Compose"
	@echo "  compose-down  - Stop services"
	@echo "  compose-logs  - Show service logs"
	@echo "  k8s-deploy    - Deploy to Kubernetes"
	@echo "  k8s-delete    - Delete from Kubernetes"
	@echo "  fmt           - Format code"
	@echo "  vet           - Vet code"
	@echo "  lint          - Lint code"
	@echo "  bench         - Benchmark analytics endpoint"
	@echo "  load-tokens   - Load tokens from file"
	@echo "  install-tools - Install development tools"
	@echo "  security-scan - Run security scan"
	@echo "  profile       - Performance profiling"
	@echo "  profile-windows - Performance profiling (Windows)"
	@echo "  test-api      - Test API with PowerShell"
	@echo "  run-windows   - Run application (Windows)"
	@echo "  help          - Show this help"

.PHONY: build clean test test-coverage test-race coverage benchmark deps tidy run run-windows docker-build docker-run docker compose-up compose-down compose-logs k8s-deploy k8s-delete fmt vet lint bench load-tokens install-tools security-scan profile profile-windows test-api help
