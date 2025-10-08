# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Lotto Numbers is a web API for generating random lottery numbers with OpenTelemetry tracing, request logging middleware, and Fly.io deployment. The application serves a static web interface and REST API endpoints.

## Architecture

### Core Structure
- **Entry point:** `cmd/lotto-numbers/main.go` - initializes HTTP server, tracing, and middleware with Gorilla Mux router
- **Number generation:** `internal/generator/` - core lottery number generation logic with shuffle-based randomization
- **HTTP handlers:** `internal/handlers/` - routes defined in `routes.go`, handler implementations in `handlers.go`
- **Static assets:** `internal/handlers/web/` - embedded filesystem for serving index.html and assets
- **Tracing:** `internal/tracing/` - Uptrace/OpenTelemetry integration with DSN-based configuration
- **Middleware:** `internal/middleware/` - logging middleware that captures request duration and strips port from IP
- **Versioning:** `internal/version/version.go` - version string injected via ldflags at build time

### Key Patterns
- **Embedded FS:** Static web assets are embedded using `//go:embed` directives and served via `fs.Sub()`
- **Tracing spans:** Most functions create tracing spans; defers ensure `span.End()` is called
- **Router setup:** All routes registered in `SetupRoutes()` which receives context for span propagation
- **Environment config:** `PORT`, `ENVIRONMENT`, and `UPTRACE_DSN` environment variables configure runtime behavior

## Development Commands

### Build
```sh
# Local build
go build -v ./...

# Build with version injection (matches CI/Docker pattern)
go build -ldflags "-s -w -X 'github.com/danstis/lotto-numbers/internal/version.Version=1.2.3'" -o lotto-numbers ./cmd/lotto-numbers
```

### Testing
```sh
# Run all tests with coverage
go test -v -coverprofile=coverage.out -covermode=count ./...

# Run tests for a specific package
go test -v ./internal/generator/

# Run a single test
go test -v -run TestGetNumbers ./internal/generator/
```

### Linting
```sh
# Run golangci-lint (configured in .golangci.yml)
golangci-lint run

# Run on specific path
golangci-lint run ./internal/handlers/
```

### Running Locally
```sh
# Set required environment variable and run
UPTRACE_DSN=your_dsn PORT=8091 go run cmd/lotto-numbers/main.go
```

## Code Quality Standards

### Error Handling with Defers
- **Never use `log.Fatal()` or `log.Fatalf()` when defers need to run** - they call `os.Exit()` immediately, bypassing defer statements
- Use `log.Panic()` or `log.Panicf()` instead to allow defers (especially `span.End()` and tracing shutdown) to execute during stack unwinding
- This is critical for proper cleanup of tracing resources and other deferred operations

### Documentation Requirements
- All packages must have package-level comments in the format: `// Package name does something.`
- Exported types must have comments in the format: `// TypeName is/does something.` (not "represents")
- Exported functions must have comments describing their purpose
- Exported variables must have comments

### Test Code
- `t.Fatalf()` is acceptable in tests (testing framework handles cleanup)
- Check all error returns, including `os.Setenv()` in tests
- Use blank identifier `_` for intentionally unused function parameters

## API Endpoints

- `GET /numbers` - Generate lottery numbers
  - Query params: `lines` (default 5), `numPerLine` (default 6), `numbersList` (comma-separated integers)
  - Returns: `{"lines": [[int]]}`
- `GET /version` - Returns version string injected at build time
- `GET /` - Serves static web interface

## Deployment

- **Docker:** Multi-stage build in `deploy/dockerfile` with version injection via `BUILD` arg
- **Fly.io:** Configured in `fly.toml`, deploy with `flyctl deploy --remote-only`
- **CI/CD:** GitHub Actions workflows handle build, test, lint, SonarCloud scanning, and deployment
- **Versioning:** GitVersion with Conventional Commits determines version numbers automatically

## External Dependencies

- **Uptrace:** Distributed tracing - requires `UPTRACE_DSN` environment variable
- **OpenTelemetry:** Instrumentation for Gorilla Mux router and runtime metrics
- **SonarCloud:** Code quality scanning in CI pipeline
- **Codecov:** Test coverage reporting
