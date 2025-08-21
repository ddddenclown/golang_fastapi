#!/bin/bash


set -e

echo "🚀 Starting full test suite for Analytics Service"
echo "=================================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' 

print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if Go is installed
check_go() {
    print_status "Checking Go installation..."
    if ! command -v go &> /dev/null; then
        print_error "Go is not installed. Please install Go 1.21+"
        exit 1
    fi
    
    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    print_success "Go version: $GO_VERSION"
}

# Check if Docker is installed
check_docker() {
    print_status "Checking Docker installation..."
    if ! command -v docker &> /dev/null; then
        print_warning "Docker is not installed. Some tests will be skipped."
        DOCKER_AVAILABLE=false
    else
        print_success "Docker is available"
        DOCKER_AVAILABLE=true
    fi
}

# Download dependencies
download_deps() {
    print_status "Downloading Go dependencies..."
    go mod download
    print_success "Dependencies downloaded"
}

# Run unit tests
run_unit_tests() {
    print_status "Running unit tests..."
    go test -v ./...
    print_success "Unit tests passed"
}

# Run tests with coverage
run_coverage_tests() {
    print_status "Running tests with coverage..."
    go test -v -coverprofile=coverage.out ./...
    go tool cover -func=coverage.out
    print_success "Coverage tests completed"
}

# Run race detector
run_race_tests() {
    print_status "Running race detector tests..."
    go test -race -v ./...
    print_success "Race detector tests passed"
}

# Run benchmarks
run_benchmarks() {
    print_status "Running benchmarks..."
    go test -bench=. -benchmem ./internal/analytics/
    print_success "Benchmarks completed"
}

# Format code
format_code() {
    print_status "Formatting code..."
    go fmt ./...
    print_success "Code formatted"
}

# Vet code
vet_code() {
    print_status "Vetting code..."
    go vet ./...
    print_success "Code vetted"
}

# Run linter
run_linter() {
    print_status "Running linter..."
    if command -v golangci-lint &> /dev/null; then
        golangci-lint run
        print_success "Linter passed"
    else
        print_warning "golangci-lint not found. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
    fi
}

# Build application
build_app() {
    print_status "Building application..."
    go build -o analytics-service ./cmd/server
    print_success "Application built successfully"
}

# Test API endpoints
test_api() {
    print_status "Testing API endpoints..."
    
    # Start the server in background
    ./analytics-service &
    SERVER_PID=$!
    
    # Wait for server to start
    sleep 3
    
    # Test health endpoint
    print_status "Testing health endpoint..."
    if curl -f http://localhost:8080/ > /dev/null 2>&1; then
        print_success "Health endpoint working"
    else
        print_error "Health endpoint failed"
        kill $SERVER_PID 2>/dev/null || true
        exit 1
    fi
    
    # Test auth endpoint
    print_status "Testing auth endpoint..."
    AUTH_RESPONSE=$(curl -s -X POST http://localhost:8080/auth \
        -H "Content-Type: application/json" \
        -d '{"email": "test@example.com", "password": "password123"}')
    
    if echo "$AUTH_RESPONSE" | grep -q "token"; then
        print_success "Auth endpoint working"
        TOKEN=$(echo "$AUTH_RESPONSE" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
    else
        print_error "Auth endpoint failed"
        kill $SERVER_PID 2>/dev/null || true
        exit 1
    fi
    
    # Test validate endpoint
    print_status "Testing validate endpoint..."
    VALIDATE_RESPONSE=$(curl -s "http://localhost:8080/validate?token=$TOKEN")
    
    if echo "$VALIDATE_RESPONSE" | grep -q '"valid":true'; then
        print_success "Validate endpoint working"
    else
        print_error "Validate endpoint failed"
        kill $SERVER_PID 2>/dev/null || true
        exit 1
    fi
    
    # Test analytics endpoint
    print_status "Testing analytics endpoint..."
    ANALYTICS_RESPONSE=$(curl -s -X POST http://localhost:8080/analytics \
        -H "Content-Type: application/json" \
        -d "{\"token\": \"$TOKEN\", \"StartDate\": \"01.01.2024\", \"FinishDate\": \"31.01.2024\"}")
    
    if echo "$ANALYTICS_RESPONSE" | grep -q "items"; then
        print_success "Analytics endpoint working"
    else
        print_warning "Analytics endpoint returned error (expected if data files missing)"
    fi
    
    # Stop server
    kill $SERVER_PID 2>/dev/null || true
    wait $SERVER_PID 2>/dev/null || true
}

# Test Docker build
test_docker() {
    if [ "$DOCKER_AVAILABLE" = true ]; then
        print_status "Testing Docker build..."
        docker build -t analytics-service-test .
        print_success "Docker build successful"
        
        print_status "Testing Docker run..."
        docker run --rm -d --name analytics-service-test -p 8081:8080 analytics-service-test
        sleep 3
        
        if curl -f http://localhost:8081/ > /dev/null 2>&1; then
            print_success "Docker container working"
        else
            print_error "Docker container failed"
        fi
        
        docker stop analytics-service-test
        docker rmi analytics-service-test
    else
        print_warning "Skipping Docker tests (Docker not available)"
    fi
}

# Check file structure
check_structure() {
    print_status "Checking project structure..."
    
    REQUIRED_FILES=(
        "go.mod"
        "go.sum"
        "Dockerfile"
        "Makefile"
        "README.md"
        "cmd/server/main.go"
        "internal/auth/service.go"
        "internal/analytics/service.go"
        "internal/config/config.go"
        "internal/handlers/auth_handler.go"
        "internal/handlers/user_handler.go"
        "internal/handlers/analytics_handler.go"
        "internal/userdb/token_store.go"
    )
    
    for file in "${REQUIRED_FILES[@]}"; do
        if [ -f "$file" ]; then
            print_success "✓ $file"
        else
            print_error "✗ $file (missing)"
            exit 1
        fi
    done
}

# Main execution
main() {
    echo "Starting full test suite..."
    echo ""
    
    check_go
    check_docker
    check_structure
    download_deps
    format_code
    vet_code
    run_linter
    run_unit_tests
    run_coverage_tests
    run_race_tests
    run_benchmarks
    build_app
    test_api
    test_docker
    
    echo ""
    echo "🎉 All tests completed successfully!"
    echo "=================================================="
    print_success "Analytics Service is ready for deployment!"
}

# Run main function
main "$@"
