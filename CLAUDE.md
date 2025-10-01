# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Budget Book (家計簿) is a full-stack personal finance management application with Go backend and Vue 3 frontend. The backend follows Clean Architecture principles with strict layer separation.

## Common Commands

### Development Setup
```bash
# Install all dependencies (backend + frontend)
make install

# Start all services with Docker
docker-compose up -d

# Start backend only (requires MySQL running)
cd backend && go run cmd/api/main.go

# Start frontend only
cd frontend && yarn dev
```

### Testing
```bash
# Run all tests (backend + frontend)
make test

# Backend tests only
cd backend && go test ./...

# Backend tests with verbose output
cd backend && go test -v ./...

# Backend tests with coverage
cd backend && go test -cover ./...

# Single package test
cd backend && go test ./usecase/
cd backend && go test ./interface/handler/

# Frontend tests (watch mode)
cd frontend && yarn test

# Frontend tests (single run)
cd frontend && yarn test:run

# Frontend tests with coverage
cd frontend && yarn test:coverage
```

### Code Quality
```bash
# Lint (backend + frontend)
make lint

# Lint with auto-fix
make lint-fix

# Format code
make format

# TypeScript type checking
make type-check
# or
cd frontend && yarn type-check
```

### Mock Generation (Backend)
```bash
# Regenerate mocks after changing repository interfaces
cd backend
mockgen -source=usecase/transaction.go -destination=mocks/repository/transaction_mock.go -package=repository
mockgen -source=usecase/category.go -destination=mocks/repository/category_mock.go -package=repository
```

## Architecture

### Backend: Clean Architecture with DI

The backend follows a strict layered architecture with dependency inversion:

```
entity/              Domain entities with business logic validation
  ↑
usecase/            Business logic + Repository interfaces (DEFINED HERE)
  ↑
infrastructure/     Repository implementations (DEPEND ON usecase interfaces)
  repository/
  ↑
interface/          HTTP handlers (depend on usecase)
  handler/
  middleware/
```

**Critical Pattern**: Repository interfaces are defined in `usecase/` packages (e.g., `TransactionRepositoryInterface` in `usecase/transaction.go`), and implementations live in `infrastructure/repository/`. This enables the dependency inversion principle.

**Dependency Injection Flow** (see `cmd/api/main.go:25-87`):
1. Create repository implementations with DB connection
2. Inject repositories into use cases (which only know about interfaces)
3. Inject use cases into handlers
4. Register handlers to Echo routes

**Testing Strategy**:
- Entity tests: Pure domain logic validation
- Usecase tests: Business logic with mocked repositories (using mockgen)
- Handler tests: Integration tests with real SQLite database
- All tests use testify for assertions

### Frontend: Vue 3 Composition API

```
views/              Page components (BudgetView, CategoryView, etc.)
  ↓
components/         Reusable UI components
  ↓
stores/ (Pinia)     State management with API integration
  ↓
services/api.ts     Axios client configuration
```

**State Management**: Pinia stores handle both state and API calls. Components interact with stores, not API directly.

**Testing**: Vitest + Vue Test Utils with jsdom. Tests live in `__tests__/` directories alongside source files.

## Key Technologies

### Backend
- Echo v4: Web framework
- GORM: ORM with MySQL (production) and SQLite (testing)
- validator.v9: Request validation
- mockgen: Type-safe mock generation for testing

### Frontend
- Vue 3 + TypeScript: Composition API with script setup
- Vuetify: Material Design components
- Pinia: State management
- Chart.js + vue-chartjs: Data visualization
- Vitest: Testing framework

## Service URLs

When running with docker-compose:
- Frontend: http://localhost:5173
- Backend API: http://localhost:8080
- MySQL: localhost:3306

## Important Patterns

### Backend Repository Interface Pattern
When creating a new feature:
1. Define entity in `entity/`
2. Define repository interface in `usecase/[feature].go`
3. Implement use case logic in same file
4. Create repository implementation in `infrastructure/repository/`
5. Create handler in `interface/handler/`
6. Wire up in `cmd/api/main.go`

### Backend Testing Pattern
- Use `mockgen` to generate mocks from repository interfaces
- Usecase tests mock repositories
- Handler tests use real SQLite database for integration testing
- See `usecase/transaction_test.go` and `interface/handler/transaction_test.go` for examples

### Frontend Component Pattern
- Use `<script setup lang="ts">` with Composition API
- Define types in `types/`
- Store business logic in Pinia stores
- Keep components focused on presentation

## Database

MySQL is used in production. SQLite is used for backend integration tests.

Connection configured via environment variables:
- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`

See `docker-compose.yml` for default values.
