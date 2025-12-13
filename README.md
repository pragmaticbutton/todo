# TODO (learning project)

A Go project for learning and experimenting with clean architecture, dependency injection, and CLI patterns. This is my sandbox for exploring Go concepts — not meant for production.

## What's in here?

- **Clean architecture** with domain, service, and storage layers.
- **CLI tool** built with Cobra for managing tasks and lists.
- **In-memory storage** with thread-safe maps (`sync.RWMutex`).
- **Service layer** handling business logic (create, update, delete, complete tasks/lists).
- **Comprehensive tests** with golden file testing and race detection.

## Quick start

Run the CLI:

```bash
make run
```

Run tests:

```bash
make test
```

Build the CLI binary:

```bash
make build-cli
```

## Project layout

- `cmd/cli/` — CLI entry point
- `cmd/todo/` — experimental sandbox entry point
- `internal/cli/` — CLI command implementations with golden tests
- `internal/domain/` — task and list domain types
- `internal/service/` — business logic layer
- `internal/storage/` — storage abstraction and memory implementation
- `internal/utils/` — utility helpers

## For curious developers

See `TODOs.md` for ideas I'm playing with (alternative storage backends, better error handling, logging, etc.).
Run `make lint` to check code quality.
