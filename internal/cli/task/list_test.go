package task_test

import (
	"bytes"
	"path/filepath"
	"testing"

	taskcmd "github.com/pragmaticbutton/todo/internal/cli/task"
	"github.com/pragmaticbutton/todo/internal/domain/task"
)

func TestListCmd_Golden(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		tasks  []task.Task
		golden string
	}{
		{
			name:   "empty",
			tasks:  []task.Task{},
			golden: filepath.Join("testdata", "list", "empty.golden"),
		},
		{
			name: "mixed",
			tasks: []task.Task{
				{ID: 1, Description: "A"},
				{ID: 2, Description: "B", Done: true},
			},
			golden: filepath.Join("testdata", "list", "mixed.golden"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := newTaskServiceWithTasks(t, tt.tasks...)
			cmd := taskcmd.NewListCmd(svc)

			var buf bytes.Buffer
			cmd.SetOut(&buf)

			if err := cmd.Execute(); err != nil {
				t.Fatalf("execute list: %v", err)
			}

			assertGolden(t, buf.String(), tt.golden)
		})
	}
}
