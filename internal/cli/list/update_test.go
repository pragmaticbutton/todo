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
	cmd.SetArgs([]string{"1", "--name", "New", "--description", "New desc"})

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
	cmd.SetArgs([]string{"1", "--name", "New", "--description", "New desc"})

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

func TestUpdateCmd_Integration(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		list    list.List
		args    []string
		checkFn func(t *testing.T, lst *list.List)
	}{
		{
			name: "update name",
			list: list.List{ID: 1, Name: "Old", Description: "Keep"},
			args: []string{"1", "--name", "New"},
			checkFn: func(t *testing.T, lst *list.List) {
				if lst.Name != "New" {
					t.Fatalf("expected name to be updated, got %s", lst.Name)
				}
				if lst.Description != "Keep" {
					t.Fatalf("expected description to remain, got %s", lst.Description)
				}
			},
		},
		{
			name: "update description",
			list: list.List{ID: 1, Name: "Old", Description: "Old desc"},
			args: []string{"1", "--description", "New desc"},
			checkFn: func(t *testing.T, lst *list.List) {
				if lst.Description != "New desc" {
					t.Fatalf("expected description to be updated, got %s", lst.Description)
				}
				if lst.Name != "Old" {
					t.Fatalf("expected name to remain, got %s", lst.Name)
				}
			},
		},
	}

	for _, tt := range tests {
		tc := tt
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			svc := newListServiceWithLists(t, tc.list)
			cmd := listcmd.NewUpdateCmd(svc)
			cmd.SetArgs(tc.args)

			if err := cmd.Execute(); err != nil {
				t.Fatalf("execute update: %v", err)
			}

			updated, err := svc.GetList(1)
			if err != nil {
				t.Fatalf("get list: %v", err)
			}

			tc.checkFn(t, updated)
		})
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
			args: []string{"abc", "--name", "New"},
		},
		{
			name: "missing list",
			args: []string{"2", "--name", "New"},
		},
		{
			name: "no fields provided",
			args: []string{"1"},
		},
	}

	for _, tt := range tests {
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
