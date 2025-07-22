# Lotto Numbers

[![Build Status](https://github.com/danstis/lotto-numbers/actions/workflows/build.yml/badge.svg)](https://github.com/danstis/lotto-numbers/actions/workflows/build.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/danstis/lotto-numbers?style=flat-square)](https://goreportcard.com/report/github.com/danstis/lotto-numbers)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/danstis/lotto-numbers)](https://pkg.go.dev/github.com/danstis/lotto-numbers)
[![Release](https://img.shields.io/github/release/danstis/lotto-numbers.svg?style=flat-square)](https://github.com/danstis/lotto-numbers/releases/latest)
[![codecov](https://codecov.io/gh/danstis/lotto-numbers/graph/badge.svg?token=csgW5w5uNs)](https://codecov.io/gh/danstis/lotto-numbers)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=danstis_lotto-numbers&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=danstis_lotto-numbers)

---

## Overview

**Lotto Numbers** is an open source web API for generating random lottery numbers, with built-in tracing, logging, and CI/CD integration. Hosted at [https://lotto-numbers.fly.dev](https://lotto-numbers.fly.dev).

## Features

- REST API for random lottery number generation
- Configurable via environment variables
- Tracing with Uptrace & OpenTelemetry
- Request logging middleware
- Static web assets served from embedded FS
- Versioning via GitVersion and CI/CD
- Docker & Fly.io deployment
- Code quality checks (SonarCloud, Codecov)

## Architecture

- **Main entry:** `cmd/lotto-numbers/main.go` (HTTP server, tracing, middleware)
- **Core logic:** `internal/generator/` (number generation)
- **Handlers:** `internal/handlers/` (routes/controllers)
- **Models:** `internal/models/` (data models)
- **Middleware:** `internal/middleware/` (logging)
- **Tracing:** `internal/tracing/` (Uptrace/OpenTelemetry)
- **Static assets:** `internal/handlers/web/`
- **Versioning:** `internal/version/version.go`

## Getting Started

### Prerequisites

- Go 1.20+
- Docker (optional, for containerized builds)
- flyctl (for deployment)

### Installation

```sh
# Clone the repo
git clone https://github.com/danstis/lotto-numbers.git
cd lotto-numbers

# Build locally
go build -v ./...

# Run tests
go test -v -coverprofile=coverage.out -covermode=count ./...

# Lint
golangci-lint run
```

### Usage

```sh
# Run the server
PORT=8080 UPTRACE_DSN=your_dsn go run cmd/lotto-numbers/main.go
```

- API endpoint: `GET /numbers`
  - Query params: `lines`, `numPerLine`, `numbersList`
  - Returns: `{ "lines": [[int]] }`
- Version endpoint: `GET /version`

## Deployment

- Docker image built and pushed to GHCR
- Deployed to Fly.io (`flyctl deploy --remote-only`)
- See `fly.toml` and `deploy/compose.yml` for config

## Contributing

Contributions are welcome! Please follow these guidelines:

- Use [Conventional Commits](https://www.conventionalcommits.org/)
- Open issues for bugs/feature requests
- Submit pull requests with clear descriptions
- Ensure all tests and lints pass

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.

## Maintainers & Contact

- [danstis](https://github.com/danstis)

## Acknowledgements

- [Uptrace](https://uptrace.dev/) for tracing
- [SonarCloud](https://sonarcloud.io/) for code quality
- [Codecov](https://codecov.io/) for coverage
- [Fly.io](https://fly.io/) for hosting

---

> _Built with ❤️ for the open source community._
