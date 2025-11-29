# TODO (learning project)

This repository contains a small TODO learning project in Go. It does not expose an HTTP/CLI interface yet; the `main.go` file simply wires services together so you can explore the domain logic.

## Current state

- In-memory storage backed by maps with `sync.RWMutex` guards to prevent data races.
- Domain types for tasks and lists, plus a service layer that supports creating, listing, updating, deleting, and completing tasks and lists.
- Monotonic ID generators in storage to avoid collisions even when records are removed.
- Utility helpers for optional values and simple calculations (e.g., percentage complete).
- No persistence layer, API, or dedicated CLI — everything runs in-process for now.

## Purpose (learning focus)

- Practice Go project layout, interfaces, and dependency injection.
- Explore simple storage abstractions and a service layer without external dependencies.
- Provide a sandbox to iterate on ideas like validation, error handling, logging, and alternative storage backends.

## Project layout (top-level)

- `main.go` — example usage / entry point for manual experimentation.
- `internal/domain` — domain types for tasks and lists.
- `internal/service` — business logic / service layer.
- `internal/storage` — storage interface and implementations (currently memory-only).
- `internal/utils` — small helpers (pointers, time helpers, etc.).

## Run locally

- Build and run with the provided Makefile:

  ```bash
  make run
  ```

- Or run directly:

  ```bash
  go run main.go
  ```

## Development helpers

- Lint: `make lint` (uses `golangci-lint`).
- Tests: `make test` (runs `go test -v -race -cover ./...`).

## Status and next steps

This is intentionally **not** production-ready. See `TODOs.md` for the roadmap (adding APIs/CLI, better error handling, alternative storage backends, etc.). The project exists for learning purposes as new Go concepts are explored.
