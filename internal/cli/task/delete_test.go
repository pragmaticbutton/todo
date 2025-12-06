package task_test

import (
	"bytes"
	"path/filepath"
	"testing"

	taskcmd "github.com/pragmaticbutton/todo/internal/cli/task"
	"github.com/pragmaticbutton/todo/internal/domain/task"
)

func TestDeleteCmd_Golden(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		task   task.Task
		args   []string
		golden string
	}{
		{
			name: "delete existing",
			task: task.Task{ID: 1, Description: "A"},
			args: []string{"1"},
			golden: filepath.Join(
				"testdata",
				"delete",
				"success.golden",
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := newTaskServiceWithTasks(t, tt.task)
			cmd := taskcmd.NewDeleteCmd(svc)

			var buf bytes.Buffer
			cmd.SetOut(&buf)
			cmd.SetArgs(tt.args)

			if err := cmd.Execute(); err != nil {
				t.Fatalf("execute delete: %v", err)
			}

			assertGolden(t, buf.String(), tt.golden)
		})
	}
}

func TestDeleteCmd_Unit(t *testing.T) {
	t.Parallel()

	svc := newTaskServiceWithTasks(t, task.Task{ID: 1, Description: "A"})
	cmd := taskcmd.NewDeleteCmd(svc)
	cmd.SetArgs([]string{"1"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute delete: %v", err)
	}

	_, err := svc.GetTask(1)
	if err == nil {
		t.Fatalf("expected task to be deleted")
	}
}

func TestDeleteCmd_InvalidID(t *testing.T) {
	t.Parallel()

	svc := newTaskServiceWithTasks(t)
	cmd := taskcmd.NewDeleteCmd(svc)
	cmd.SetArgs([]string{"abc"})

	if err := cmd.Execute(); err == nil {
		t.Fatalf("expected error for invalid id")
	}
}
