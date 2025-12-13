package list_test

import (
	"bytes"
	"path/filepath"
	"testing"

	listcmd "github.com/pragmaticbutton/todo/internal/cli/list"
	listdomain "github.com/pragmaticbutton/todo/internal/domain/list"
)

func TestAddCmd_Golden(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		args   []string
		golden string
	}{
		{
			name:   "simple add",
			args:   []string{"Groceries"},
			golden: filepath.Join("testdata", "add", "simple_add.golden"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := newListServiceWithLists(t)
			cmd := listcmd.NewAddCmd(svc)

			var buf bytes.Buffer
			cmd.SetOut(&buf)
			cmd.SetArgs(tt.args)

			if err := cmd.Execute(); err != nil {
				t.Fatalf("execute add: %v", err)
			}

			assertGolden(t, buf.String(), tt.golden)
		})
	}
}

func TestAddCmd_Integration(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		args    []string
		checkFn func(t *testing.T, lst *listdomain.List)
	}{
		{
			name: "add list with description",
			args: []string{"Groceries", "--description", "Weekly shopping"},
			checkFn: func(t *testing.T, lst *listdomain.List) {
				if lst.Name != "Groceries" {
					t.Fatalf("unexpected name: %s", lst.Name)
				}
				if lst.Description != "Weekly shopping" {
					t.Fatalf("unexpected description: %s", lst.Description)
				}
			},
		},
		{
			name: "add list without description",
			args: []string{"Work"},
			checkFn: func(t *testing.T, lst *listdomain.List) {
				if lst.Name != "Work" {
					t.Fatalf("unexpected name: %s", lst.Name)
				}
				if lst.Description != "" {
					t.Fatalf("expected empty description, got %s", lst.Description)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := newListServiceWithLists(t)
			cmd := listcmd.NewAddCmd(svc)
			cmd.SetArgs(tt.args)

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

			tt.checkFn(t, &lists[0])
		})
	}
}

func TestAddCmd_Errors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		args []string
	}{
		{
			name: "missing name",
			args: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := newListServiceWithLists(t)
			cmd := listcmd.NewAddCmd(svc)
			cmd.SetArgs(tt.args)

			if err := cmd.Execute(); err == nil {
				t.Fatalf("expected error for case %q", tt.name)
			}
		})
	}
}
