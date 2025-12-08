package list_test

import (
	"bytes"
	"path/filepath"
	"testing"

	listcmd "github.com/pragmaticbutton/todo/internal/cli/list"
)

func TestAddCmd_Golden(t *testing.T) {
	t.Parallel()

	svc := newListServiceWithLists(t)
	cmd := listcmd.NewAddCmd(svc)

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"Groceries", "Weekly shopping"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute add: %v", err)
	}

	goldenPath := filepath.Join("testdata", "add", "simple_add.golden")
	assertGolden(t, buf.String(), goldenPath)
}

func TestAddCmd_Unit(t *testing.T) {
	t.Parallel()

	svc := newListServiceWithLists(t)
	cmd := listcmd.NewAddCmd(svc)
	cmd.SetArgs([]string{"Groceries", "Weekly shopping"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute add: %v", err)
	}

	lists, err := svc.ListLists()
	if err != nil {
		t.Fatalf("list lists: %v", err)
	}
	if len(lists) != 1 {
		t.Fatalf("expected 1 list, got %d", len(lists))
	}
	if lists[0].Name != "Groceries" || lists[0].Description != "Weekly shopping" {
		t.Fatalf("unexpected list data: %+v", lists[0])
	}
}
