.PHONY: fmt lint test build clean

# Format Go code
fmt:
	gofmt -s -w .
	goimports -w .

# Run linter
lint:
	golangci-lint run

# Fix linting issues where possible
lint-fix:
	golangci-lint run --fix

# Run tests
test:
	go test -v ./...

# Build the application
build:
	go build -o bin/main ./cmd/main.go

# Clean build artifacts
clean:
	rm -rf bin/

# Run all checks
check: fmt lint test

# Install development tools
install-tools:
	go install golang.org/x/tools/cmd/goimports@latest
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.55.2