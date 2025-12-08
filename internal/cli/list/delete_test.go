package list_test

import (
	"bytes"
	"path/filepath"
	"testing"

	listcmd "github.com/pragmaticbutton/todo/internal/cli/list"
	"github.com/pragmaticbutton/todo/internal/domain/list"
)

func TestDeleteCmd_Golden(t *testing.T) {
	t.Parallel()

	svc := newListServiceWithLists(t, list.List{ID: 1, Name: "Groceries", Description: "Weekly"})
	cmd := listcmd.NewDeleteCmd(svc)

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"1"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute delete: %v", err)
	}

	goldenPath := filepath.Join("testdata", "delete", "success.golden")
	assertGolden(t, buf.String(), goldenPath)
}

func TestDeleteCmd_Unit(t *testing.T) {
	t.Parallel()

	svc := newListServiceWithLists(t, list.List{ID: 1, Name: "Groceries"})
	cmd := listcmd.NewDeleteCmd(svc)
	cmd.SetArgs([]string{"1"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute delete: %v", err)
	}

	if _, err := svc.GetList(1); err == nil {
		t.Fatalf("expected list to be deleted")
	}
}

func TestDeleteCmd_InvalidID(t *testing.T) {
	t.Parallel()

	svc := newListServiceWithLists(t)
	cmd := listcmd.NewDeleteCmd(svc)
	cmd.SetArgs([]string{"abc"})

	if err := cmd.Execute(); err == nil {
		t.Fatalf("expected error for invalid id")
	}
}
