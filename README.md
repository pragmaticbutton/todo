# TODO (learning project)

This is a small in-memory TODO application created for learning and experimentation with Go.

Purpose

- Practice Go project layout, interfaces and dependency injection.
- Explore simple in-memory storage, services and basic domain modeling.

Features

- Add, list, get, update and delete tasks and lists.
- In-memory storage implementation (no external DB).
- Simple service layer and domain types.

Project layout (top-level)

- `main.go` — example usage / entry point
- `internal/domain` — domain types for tasks and lists
- `internal/service` — business logic / service layer
- `internal/storage` — storage interface and implementations (memory)
- `internal/utils` — small helpers

Run

- Build and run with the provided Makefile:

  make run

- A VS Code launch task (`.vscode/launch.json`) is included that runs the `make run` task.

Notes

- This repository is intended for learning and experimentation only; it is not production-ready.
- Concurrency and persistence are intentionally simple; consider adding synchronization or a real database for production use.
