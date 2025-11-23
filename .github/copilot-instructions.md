# Copilot Instructions for TODO Project

A learning Go project exploring layered architecture, dependency injection, and in-memory storage patterns.

## Architecture Overview

**Layered Design**: The project follows a clean architecture with clear separation of concerns:

- **Domain** (`internal/domain/{task,list}/`): Pure data types with minimal logic. `Task` and `List` are PODs with ID, timestamps, and constructor helpers.
- **Service** (`internal/service/`): Business logic layer. `ListService` and `TaskService` orchestrate domain operations via the storage interface.
- **Storage** (`internal/storage/`): Abstraction layer defining the `Storage` interface. Current implementation is in-memory (`memory.go`), intentionally thread-unsafe for learning.
- **Main** (`main.go`): Dependency injection entry point—instantiates storage, then services.

**Key Insight**: Storage is injected into services, services are never aware of storage implementation details. Adding a database implementation only requires a new package satisfying the `Storage` interface.

## Testing Patterns

Uses `testify` (assert/require) extensively. Test pattern in `internal/storage/memory/memory_test.go`:

- **Subtests**: `t.Run("scenario", func(t *testing.T) {...})` for test organization
- **Setup/Reset**: `resetForTest()` resets the singleton memory storage between tests (see line ~440 in memory_test.go for reset helper)
- **Fixtures**: Fixed timestamps (`time.Date(...)`) for deterministic testing
- **Assertions**: `require.NoError()` for early exit on error, `assert.Equal()` for comparisons

Run tests with: `make test` (includes `-race` flag for concurrency detection)

## Project Conventions

**Optional Fields**: Use pointer types (`*uint32`, `*string`) for optional values. Helper function `utils.Ptr()` converts values to pointers.

**Error Handling**: Errors are descriptive, including context (e.g., `"task with id %d not found"`). Services return errors directly without custom error types.

**Timestamps**: Domain types include `Created` and `Updated` fields. Services update `Updated` on mutations (e.g., `UpdateList` sets `time.Now()`). Both are `time.Time` (not pointers)—TODO comments suggest this pattern needs discussion for nil semantics.

**ID Generation**: Services call `storage.NextTaskID()`/`storage.NextListID()` to get next IDs. Memory implementation uses `len(map) + 1` (naive, non-persistent approach).

**Service Input Types**: Services define input DTOs (`AddListInput`, `UpdateListInput`) rather than accepting domain types directly. This decouples request contracts from domain entities.

## Build & Run

- **Run**: `make run` — builds and executes `main.go`
- **Test**: `make test` — runs all tests with `-race` flag and coverage
- **Lint**: `make lint` — runs golangci-lint (requires external installation)

VS Code launch task (`.vscode/launch.json`) invokes `make run`.

## Known TODOs (Discoverable in Code)

- Thread safety: Memory storage uses singleton pattern without mutex protection for maps
- Time.Time nil semantics: Both `Task.Updated` and `List.Updated` initialized to zero time; consider nullable approach
- Storage interface design: Monolithic interface mixes task/list methods—consider splitting
- Service layer organization: Both `ListService` and `TaskService` in one file

## File Reference Map

| Component                 | Key Files                                                                                                    |
| ------------------------- | ------------------------------------------------------------------------------------------------------------ |
| Domain Types              | `internal/domain/task/task.go`, `internal/domain/list/list.go`                                               |
| Service Layer             | `internal/service/list.go`, `internal/service/task.go`                                                       |
| Storage Interface & Tests | `internal/storage/storage.go`, `internal/storage/memory/memory.go`, `internal/storage/memory/memory_test.go` |
| Utilities                 | `internal/utils/utils.go`                                                                                    |
| Entry Point               | `main.go`                                                                                                    |
