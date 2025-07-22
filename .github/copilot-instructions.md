# Copilot Instructions for lotto-numbers

## Project Overview

-   **Purpose:** Generates random lottery numbers via a web API, with tracing and logging, deployed on Fly.io.
-   **Architecture:**
    -   Main entry: `cmd/lotto-numbers/main.go` sets up HTTP server, tracing, and middleware.
    -   Core logic: `internal/generator/` (number generation), `internal/handlers/` (HTTP routes/controllers), `internal/models/` (data models).
    -   Middleware: `internal/middleware/` for logging requests.
    -   Tracing: `internal/tracing/` integrates Uptrace and OpenTelemetry.
    -   Static web assets: `internal/handlers/web/` (served via embedded FS).
    -   Versioning: `internal/version/version.go` and GitVersion workflow.

## Key Workflows

-   **Build:**
    -   Local: `go build -v ./...`
    -   Docker: See `deploy/dockerfile` for multi-stage build and version injection.
    -   CI: GitHub Actions (`.github/workflows/build.yml`) runs build, tests, lint, and SonarCloud scan.
-   **Test:**
    -   Local: `go test -v -coverprofile=coverage.out -covermode=count ./...`
    -   CI: Coverage uploaded to Codecov.
-   **Lint:**
    -   Local: `golangci-lint run`
    -   CI: `golangci-lint-action` in workflow.
-   **Release/Deploy:**
    -   Versioning via GitVersion and Conventional Commits.
    -   Docker image built and pushed to GHCR, deployed to Fly.io (`flyctl deploy --remote-only`).
    -   See `fly.toml` and `deploy/compose.yml` for config.

## Patterns & Conventions

-   **Routing:** All HTTP routes registered in `internal/handlers/routes.go` via Gorilla Mux. Static files served from embedded FS.
-   **API:** `/numbers` endpoint returns JSON `{ "lines": [[int]] }`. Query params: `lines`, `numPerLine`, `numbersList`.
-   **Tracing:** Uptrace DSN required via `UPTRACE_DSN` env var. Tracing setup in `internal/tracing/tracing.go`.
-   **Logging:** All requests logged with IP, method, path, and duration.
-   **Version:** Injected at build time via `-ldflags` and surfaced at `/version` endpoint.
-   **Environment:** Configurable via env vars (`PORT`, `ENVIRONMENT`, `UPTRACE_DSN`). Defaults set in Dockerfile and Fly.io config.

## External Integrations

-   **Uptrace:** Distributed tracing (see `internal/tracing/`).
-   **SonarCloud:** Code quality scan in CI.
-   **Codecov:** Test coverage reporting.
-   **GitHub Container Registry (GHCR):** Docker image hosting.
-   **Fly.io:** App hosting and deployment.

## Examples

-   To add a new API route, update `internal/handlers/routes.go` and implement handler in `internal/handlers/handlers.go`.
-   To change lottery logic, edit `internal/generator/generator.go`.
-   To update version, use Conventional Commit and let CI workflows handle tagging and injection.
