package task_test

import (
	"flag"
	"os"
	"testing"

	"github.com/pragmaticbutton/todo/internal/domain/task"
	"github.com/pragmaticbutton/todo/internal/service"
	"github.com/pragmaticbutton/todo/internal/storage/memory"
)

var update = flag.Bool("update", false, "update golden files")

func assertGolden(t *testing.T, actual string, goldenPath string) {
	t.Helper()

	if *update {
		if err := os.WriteFile(goldenPath, []byte(actual), 0o600); err != nil {
			t.Fatalf("write golden file %s: %v", goldenPath, err)
		}
	}

	expected, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("read golden file %s: %v", goldenPath, err)
	}

	if string(expected) != actual {
		t.Fatalf("golden mismatch\nexpected:\n%s\ngot:\n%s", string(expected), actual)
	}
}

func newTaskServiceWithTasks(t *testing.T, tasks ...task.Task) *service.TaskService {
	t.Helper()

	store := memory.New()
	for _, tk := range tasks {
		if err := store.AddTask(&tk); err != nil {
			t.Fatalf("seed task: %v", err)
		}
	}

	return service.NewTaskService(store, store)
}
