package list_test

import (
	"bytes"
	"path/filepath"
	"testing"

	listcmd "github.com/pragmaticbutton/todo/internal/cli/list"
	"github.com/pragmaticbutton/todo/internal/domain/list"
	"github.com/pragmaticbutton/todo/internal/domain/task"
	"github.com/pragmaticbutton/todo/internal/utils"
)

func TestListTasksCmd_Golden(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		tasks  []task.Task
		golden string
	}{
		{
			name:   "empty",
			tasks:  []task.Task{},
			golden: filepath.Join("testdata", "tasks", "empty.golden"),
		},
		{
			name: "mixed",
			tasks: []task.Task{
				{ID: 2, ListID: utils.Ptr(uint32(1)), Description: "B", Done: true},
				{ID: 1, ListID: utils.Ptr(uint32(1)), Description: "A"},
				{ID: 3, ListID: utils.Ptr(uint32(2)), Description: "Other list"},
			},
			golden: filepath.Join("testdata", "tasks", "mixed.golden"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := newListServiceWithData(t, []list.List{{ID: 1, Name: "Groceries"}}, tt.tasks)
			cmd := listcmd.NewListTasksCmd(svc)

			var buf bytes.Buffer
			cmd.SetOut(&buf)
			cmd.SetArgs([]string{"1"})

			if err := cmd.Execute(); err != nil {
				t.Fatalf("execute list tasks: %v", err)
			}

			assertGolden(t, buf.String(), tt.golden)
		})
	}
}

func TestListTasksCmd_InvalidID(t *testing.T) {
	t.Parallel()

	svc := newListServiceWithLists(t, list.List{ID: 1, Name: "Groceries"})
	cmd := listcmd.NewListTasksCmd(svc)
	cmd.SetArgs([]string{"abc"})

	if err := cmd.Execute(); err == nil {
		t.Fatalf("expected error for invalid id")
	}
}
