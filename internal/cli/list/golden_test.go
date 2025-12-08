package list_test

import (
	"flag"
	"os"
	"testing"

	"github.com/pragmaticbutton/todo/internal/domain/list"
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

func newListServiceWithLists(t *testing.T, lists ...list.List) *service.ListService {
	t.Helper()
	return newListServiceWithData(t, lists, nil)
}

func newListServiceWithData(t *testing.T, lists []list.List, tasks []task.Task) *service.ListService {
	t.Helper()

	store := memory.New()
	for _, l := range lists {
		if err := store.AddList(&l); err != nil {
			t.Fatalf("seed list: %v", err)
		}
	}
	for _, tk := range tasks {
		if err := store.AddTask(&tk); err != nil {
			t.Fatalf("seed task: %v", err)
		}
	}

	return service.NewListService(store, store)
}
