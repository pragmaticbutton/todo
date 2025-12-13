package list_test

import (
	"bytes"
	"path/filepath"
	"testing"

	listcmd "github.com/pragmaticbutton/todo/internal/cli/list"
	"github.com/pragmaticbutton/todo/internal/domain/list"
	"github.com/pragmaticbutton/todo/internal/service"
)

func TestDeleteCmd_Golden(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		lists []list.List
		args  []string
	}{
		{
			name:  "delete existing",
			lists: []list.List{{ID: 1, Name: "Groceries", Description: "Weekly"}},
			args:  []string{"1"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := newListServiceWithLists(t, tt.lists...)
			cmd := listcmd.NewDeleteCmd(svc)

			var buf bytes.Buffer
			cmd.SetOut(&buf)
			cmd.SetArgs(tt.args)

			if err := cmd.Execute(); err != nil {
				t.Fatalf("execute delete: %v", err)
			}

			goldenPath := filepath.Join("testdata", "delete", "success.golden")
			assertGolden(t, buf.String(), goldenPath)
		})
	}
}

func TestDeleteCmd_Integration(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		lists   []list.List
		args    []string
		checkFn func(t *testing.T, svc *service.ListService)
	}{
		{
			name:  "delete single list",
			lists: []list.List{{ID: 1, Name: "Groceries"}},
			args:  []string{"1"},
			checkFn: func(t *testing.T, svc *service.ListService) {
				if _, err := svc.GetList(1); err == nil {
					t.Fatalf("expected list to be deleted")
				}
			},
		},
		{
			name:  "delete only targeted list",
			lists: []list.List{{ID: 1, Name: "Groceries"}, {ID: 2, Name: "Work"}},
			args:  []string{"1"},
			checkFn: func(t *testing.T, svc *service.ListService) {
				if _, err := svc.GetList(1); err == nil {
					t.Fatalf("expected list to be deleted")
				}
				remaining, err := svc.GetList(2)
				if err != nil {
					t.Fatalf("expected other list to remain: %v", err)
				}
				if remaining.Name != "Work" {
					t.Fatalf("unexpected remaining list: %+v", remaining)
				}
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := newListServiceWithLists(t, tt.lists...)
			cmd := listcmd.NewDeleteCmd(svc)
			cmd.SetArgs(tt.args)

			if err := cmd.Execute(); err != nil {
				t.Fatalf("execute delete: %v", err)
			}

			tt.checkFn(t, svc)
		})
	}
}

func TestDeleteCmd_Errors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		lists []list.List
		args  []string
	}{
		{
			name: "invalid id",
			args: []string{"abc"},
		},
		{
			name:  "missing list",
			lists: []list.List{{ID: 1, Name: "Groceries"}},
			args:  []string{"2"},
		},
		{
			name: "missing argument",
			args: []string{},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := newListServiceWithLists(t, tt.lists...)
			cmd := listcmd.NewDeleteCmd(svc)
			cmd.SetArgs(tt.args)

			if err := cmd.Execute(); err == nil {
				t.Fatalf("expected error for case %q", tt.name)
			}
		})
	}
}
