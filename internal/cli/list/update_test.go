package list_test

import (
	"bytes"
	"path/filepath"
	"testing"

	listcmd "github.com/pragmaticbutton/todo/internal/cli/list"
	"github.com/pragmaticbutton/todo/internal/domain/list"
)

func TestUpdateCmd_Golden(t *testing.T) {
	t.Parallel()

	svc := newListServiceWithLists(t, list.List{ID: 1, Name: "Old", Description: "Old desc"})
	cmd := listcmd.NewUpdateCmd(svc)

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"1", "New", "New desc"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute update: %v", err)
	}

	goldenPath := filepath.Join("testdata", "update", "updated.golden")
	assertGolden(t, buf.String(), goldenPath)
}

func TestUpdateCmd_Unit(t *testing.T) {
	t.Parallel()

	svc := newListServiceWithLists(t, list.List{ID: 1, Name: "Old", Description: "Old desc"})
	cmd := listcmd.NewUpdateCmd(svc)
	cmd.SetArgs([]string{"1", "New", "New desc"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute update: %v", err)
	}

	updated, err := svc.GetList(1)
	if err != nil {
		t.Fatalf("get list: %v", err)
	}
	if updated.Name != "New" || updated.Description != "New desc" {
		t.Fatalf("expected list to be updated, got %+v", updated)
	}
}

func TestUpdateCmd_Errors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		args []string
	}{
		{
			name: "invalid id",
			args: []string{"abc", "New", "New desc"},
		},
		{
			name: "missing list",
			args: []string{"2", "New", "New desc"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := newListServiceWithLists(t, list.List{ID: 1, Name: "Old", Description: "Old desc"})
			cmd := listcmd.NewUpdateCmd(svc)
			cmd.SetArgs(tt.args)

			if err := cmd.Execute(); err == nil {
				t.Fatalf("expected error for case %q", tt.name)
			}
		})
	}
}
