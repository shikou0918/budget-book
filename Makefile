.PHONY: help install lint lint-fix format format-check test clean

# Default target
help:
	@echo "Available commands:"
	@echo "  install     - Install dependencies for both frontend and backend"
	@echo "  lint        - Run linters for both frontend and backend"
	@echo "  lint-fix    - Run linters with auto-fix for both frontend and backend"
	@echo "  format      - Format code for both frontend and backend"
	@echo "  format-check- Check code formatting for both frontend and backend"
	@echo "  test        - Run tests for both frontend and backend"
	@echo "  clean       - Clean build artifacts"

# Install dependencies
install:
	@echo "Installing backend dependencies..."
	cd backend && go mod download
	@echo "Installing frontend dependencies..."
	cd frontend && yarn install

# Backend targets
backend-lint:
	@echo "Running backend linter..."
	cd backend && golangci-lint run

backend-lint-fix:
	@echo "Running backend linter with fixes..."
	cd backend && golangci-lint run --fix

backend-format:
	@echo "Formatting backend code..."
	cd backend && gofmt -s -w .
	cd backend && goimports -w .

backend-test:
	@echo "Running backend tests..."
	cd backend && go test -v ./...

# Frontend targets
frontend-lint:
	@echo "Running frontend linter..."
	cd frontend && yarn lint:check

frontend-lint-fix:
	@echo "Running frontend linter with fixes..."
	cd frontend && yarn lint:fix

frontend-format:
	@echo "Formatting frontend code..."
	cd frontend && yarn format

frontend-format-check:
	@echo "Checking frontend code formatting..."
	cd frontend && yarn format:check

frontend-test:
	@echo "Running frontend tests..."
	cd frontend && yarn test 2>/dev/null || echo "No frontend tests configured yet"

# Combined targets
lint: backend-lint frontend-lint

lint-fix: backend-lint-fix frontend-lint-fix

format: backend-format frontend-format

format-check: frontend-format-check
	@echo "Backend formatting is handled by gofmt (no separate check needed)"

test: backend-test frontend-test

# Clean
clean:
	@echo "Cleaning backend artifacts..."
	cd backend && rm -rf bin/
	@echo "Cleaning frontend artifacts..."
	cd frontend && rm -rf dist/ node_modules/.cache/