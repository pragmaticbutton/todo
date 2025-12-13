# Copilot Instructions for TODO Project

A learning Go project exploring clean layered architecture, dependency injection, thread-safe in-memory storage, and CLI patterns with comprehensive testing.

## Architecture Overview

### Layered Design

The project strictly separates concerns across four layers:

1. **Domain Layer** (`internal/domain/{task,list}/`)

   - Pure data structures: `Task` and `List` with ID, timestamps, and validation enums
   - `Task` has: ID, Description, Done (bool), Priority (enum), ListID (optional), Created, Updated (pointer to time)
   - `List` has: ID, Name, Description, Created, Updated (pointer to time)
   - Constructor functions: `task.New(id, desc, priority, listID)` and `list.New(id, name, desc)`
   - `Priority` enum: `PriorityLow`, `PriorityMedium`, `PriorityHigh` with `ParsePriority(s string)` parser
   - **Rule**: Domain types contain NO business logic, only data and basic constructors

2. **Service Layer** (`internal/service/{list,task}.go`)

   - Orchestrates business operations using domain and storage layers
   - `ListService` and `TaskService` both injected with storage interface (not implementations)
   - Input types are DTOs: `AddListInput{Name, Description}`, `UpdateListInput{Name*, Description*}`, `AddTaskInput{Description, ListID*, Priority*}`, `UpdateTaskInput{Description*, Done*, Priority*}`
   - Each mutation calls `storage.NextTaskID()` or `storage.NextListID()` to get unique IDs
   - Mutations set `time.Now()` on `Updated` field via `utils.Ptr()`
   - **Pattern**: Service methods never access storage directly; always go through interface contract
   - **All public methods return errors directly** (no custom error types, simple string errors)

3. **Storage Interface** (`internal/storage/storage.go`)

   - Two focused interfaces: `TaskStorage` and `ListStorage` (composed into `Storage`)
   - Task methods: `NextTaskID()`, `AddTask()`, `ListTasks()`, `GetTask()`, `DeleteTask()`, `UpdateTask()`, `SearchTasks(listID *uint32)`
   - List methods: `NextListID()`, `AddList()`, `ListLists()`, `GetList()`, `DeleteList()`, `UpdateList()`
   - **Critical**: Pointers used throughout (`*task.Task`, `*list.List`) to allow mutation tracking

4. **Memory Storage** (`internal/storage/memory/memory.go`)

   - Implements full `Storage` interface with maps: `tasks map[uint32]*task.Task`, `lists map[uint32]*list.List`
   - Thread-safe with `sync.RWMutex` (acquire lock in Add/Update/Delete, release with defer)
   - Monotonic ID counters: `nextTaskID`, `nextListID` (prevents collisions even after deletions)
   - **Search pattern**: `SearchTasks(listID *uint32)` filters by ListID if not nil
   - **Error format**: Always `fmt.Errorf("entity with id %d not found", id)`

5. **CLI Layer** (`internal/cli/`)
   - Built with Cobra framework for command structure
   - Root command: `NewRootCmd(taskService, listService)` wires services
   - Task subcommands: `list`, `add`, `get`, `delete`, `update`, `complete`, `reopen` (in `internal/cli/task/`)
   - List subcommands: `list`, `add`, `get`, `delete`, `update`, `tasks` (in `internal/cli/list/`)
   - Each command uses Cobra flags via `cmd.Flags().StringVarP()`, `IntVarP()`, etc.
   - Output uses `fmt.Fprintf(cmd.OutOrStdout(), ...)` (testable)
   - **Arg pattern**: Positional args for primary input (e.g., task name), flags for optional (e.g., `-d` for description)

### Dependency Injection Flow

```
main.go
  ↓
storage := memory.New()
  ↓
services := NewListService/NewTaskService(storage, storage)
  ↓
cli.NewRootCmd(taskService, listService)
```

Services never know about concrete storage; only interface contract. New storage backends (DB, Redis) only require implementing the `Storage` interface.

## Code Patterns & Conventions

### Optional Fields

- **In domain types**: Use `*T` pointers for optional fields (e.g., `ListID *uint32`, `Updated *time.Time`)
- **In service inputs**: Use `*T` for optional fields (e.g., `UpdateTaskInput{Priority *Priority}`)
- **Helper**: `utils.Ptr[T any](v T) *T` converts any value to pointer
- **Nil check pattern**: `if input.Name != nil { field = *input.Name }`

### Error Handling

- Services return `error` directly (no custom error types)
- Format: descriptive with context, e.g., `fmt.Errorf("task with id %d not found", id)`
- CLI commands use `RunE` which propagates errors to Cobra
- Test assertions: Use `require.NoError(t, err)` to fail fast; `assert.Equal(t, expected, actual)` for comparisons

### Timestamps

- Both `Created` and `Updated` are `time.Time` fields (not pointers)
- `Created` is set once in constructor via `time.Now()`
- `Updated` is `*time.Time` (pointer), nil until first mutation, set via `utils.Ptr(time.Now())`
- This allows distinguishing "never modified" from "modified"

### ID Generation

- Services call `storage.NextTaskID()` / `storage.NextListID()` BEFORE creating domain object
- Memory implementation increments monotonic counter to prevent reuse even after deletion
- Pass ID to constructor: `task.New(storage.NextTaskID(), ...)`

### Testing Patterns

- Use `t.Parallel()` for concurrent test execution
- Use `t.Run("scenario", func(t *testing.T) {...})` for subtests
- Setup pattern: `svc := newListServiceWithLists(t)` (test helper that creates fresh storage + service)
- Golden files: Store expected CLI output in `testdata/{cmd}/`, run tests with `-update` to regenerate
- Testify assertions:
  - `require.NoError(t, err)` exits test on error
  - `assert.Equal(t, expected, actual)` continues test on mismatch
  - `require.NotNil(t, val)` fails if nil

### Command Implementation Pattern (Cobra)

```go
func NewAddCmd(svc *service.Service) *cobra.Command {
    var flagValue string
    cmd := &cobra.Command{
        Use: "add <positional>",
        Args: cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            result := svc.Add(service.AddInput{...})
            fmt.Fprintf(cmd.OutOrStdout(), "Success message\n")
            return nil
        },
    }
    cmd.Flags().StringVarP(&flagValue, "flag-name", "f", "", "description")
    return cmd
}
```

## File Organization & Quick Reference

| Purpose           | Files                                                        | Key Details                                          |
| ----------------- | ------------------------------------------------------------ | ---------------------------------------------------- |
| **Domain - Task** | `internal/domain/task/task.go`                               | Priority enum, Task struct, New(), ParsePriority()   |
| **Domain - List** | `internal/domain/list/list.go`                               | List struct, New()                                   |
| **Service Layer** | `internal/service/list.go`, `internal/service/task.go`       | ListService, TaskService, input DTOs                 |
| **Storage Iface** | `internal/storage/storage.go`                                | TaskStorage, ListStorage interfaces                  |
| **Memory Impl**   | `internal/storage/memory/memory.go`                          | Maps with RWMutex, monotonic ID counters             |
| **Memory Tests**  | `internal/storage/memory/memory_test.go`                     | Subtests, fixtures, assertions                       |
| **CLI Root**      | `internal/cli/root.go`                                       | NewRootCmd(), registers task & list subcommands      |
| **Task Commands** | `internal/cli/task/{add,list,get,delete,update,complete}.go` | Each has `_test.go` with golden files in `testdata/` |
| **List Commands** | `internal/cli/list/{add,list,get,delete,update,tasks}.go`    | Similar structure to task commands                   |
| **Utilities**     | `internal/utils/utils.go`                                    | `Ptr[T]()` generic pointer converter                 |
| **Entry Points**  | `cmd/cli/main.go`, `cmd/todo/main.go`                        | CLI wiring and experimental sandbox                  |

## Build, Test, Run

- **Run CLI**: `make run` → executes `cmd/cli/main.go`
- **Run Tests**: `make test` → `go test -v -race -cover ./...`
  - `-race` detects data race conditions
  - `-cover` shows coverage percentage
- **Build CLI Binary**: `make build-cli` → outputs to `bin/todo`
- **Update Golden Files**: `go test -v -race -cover ./internal/cli/... -args -update`
- **Lint**: `make lint` → golangci-lint (requires installation)

## Code Quality Guidelines

### When Adding Features

1. **Start with domain types** if introducing new concepts
2. **Add to storage interface** before implementation
3. **Implement in memory storage** with proper locking
4. **Write tests** in same package (with `_test.go` suffix)
5. **Add CLI command** if user-facing feature
6. **Write golden tests** for CLI output consistency

### Locking Pattern (Critical!)

```go
// Read operations: RLock/RUnlock
m.mu.RLock()
val, ok := m.tasks[id]
m.mu.RUnlock()

// Write operations: Lock/Unlock with defer
m.mu.Lock()
defer m.mu.Unlock()
m.tasks[id] = val
```

**Never hold locks across function calls** to prevent deadlocks.

### Naming Conventions

- **Interfaces**: PascalCase ending in "Storage" or action verb (e.g., `TaskStorage`, `ListService`)
- **Concrete types**: lowercase package names (e.g., `memory.New()`)
- **Methods**: PascalCase, verbs first (e.g., `AddTask`, `ListLists`, `GetTask`)
- **Flags**: kebab-case (e.g., `--description`, `-d`)
- **Packages**: lowercase single word (domain, service, storage, cli)

## Known Design Patterns & TODOs

| Pattern/Issue         | Status      | Details                                                                                         |
| --------------------- | ----------- | ----------------------------------------------------------------------------------------------- |
| Thread safety         | ✅ Complete | RWMutex protects all map access; lock/unlock pattern in every method                            |
| ID generation         | ✅ Complete | Monotonic counters prevent collisions; incremented even on deletions                            |
| Optional field syntax | ✅ Complete | Pointers + nil checks throughout; `utils.Ptr()` helper available                                |
| CLI structure         | ✅ Complete | Cobra with tests; golden files for output regression testing                                    |
| Error handling        | 🔄 WIP      | Currently simple fmt.Errorf; consider custom error types for better context                     |
| Logging               | ❌ TODO     | No logging layer; consider adding slog for diagnostics                                          |
| Storage backends      | ❌ TODO     | Only memory impl; design allows DB/Redis/file implementations by satisfying `Storage` interface |
| Validation            | ❌ TODO     | No input validation (e.g., empty strings, negative priority); consider validation layer         |
| API layer             | ❌ TODO     | Currently CLI-only; REST/gRPC could be added without changing services                          |

## Copilot Usage Tips

- **Ask for tests**: When writing new features, request test coverage first
- **Ask for error handling**: Specify error cases to avoid (nil dereferences, race conditions)
- **Reference patterns**: Say "follow the pattern in `TaskService.UpdateTask`" for consistency
- **Clarify layer**: Specify whether change is domain/service/storage/CLI for proper placement
- **Linting**: Request golangci-lint compliance when writing code
- **Thread safety**: Always mention data races and mutex usage for storage layer changes
