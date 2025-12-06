package task_test

import (
	"bytes"
	"path/filepath"
	"testing"

	taskcmd "github.com/pragmaticbutton/todo/internal/cli/task"
	"github.com/pragmaticbutton/todo/internal/domain/task"
)

func TestCompleteCmd_Golden(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		task   task.Task
		args   []string
		golden string
	}{
		{
			name:   "mark done",
			task:   task.Task{ID: 1, Description: "Do it"},
			args:   []string{"1"},
			golden: filepath.Join("testdata", "complete", "completed.golden"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := newTaskServiceWithTasks(t, tt.task)
			cmd := taskcmd.NewCompleteCmd(svc)

			var buf bytes.Buffer
			cmd.SetOut(&buf)
			cmd.SetArgs(tt.args)

			if err := cmd.Execute(); err != nil {
				t.Fatalf("execute complete: %v", err)
			}

			assertGolden(t, buf.String(), tt.golden)
		})
	}
}

func TestCompleteCmd_Unit(t *testing.T) {
	t.Parallel()

	svc := newTaskServiceWithTasks(t, task.Task{ID: 1, Description: "Do it"})
	cmd := taskcmd.NewCompleteCmd(svc)
	cmd.SetArgs([]string{"1"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute complete: %v", err)
	}

	tk, err := svc.GetTask(1)
	if err != nil {
		t.Fatalf("get task: %v", err)
	}
	if !tk.Done {
		t.Fatalf("expected task to be marked done")
	}
}

func TestCompleteCmd_Errors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		args []string
	}{
		{
			name: "invalid id",
			args: []string{"abc"},
		},
		{
			name: "missing task",
			args: []string{"2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := newTaskServiceWithTasks(t, task.Task{ID: 1, Description: "Do it"})
			cmd := taskcmd.NewCompleteCmd(svc)
			cmd.SetArgs(tt.args)

			if err := cmd.Execute(); err == nil {
				t.Fatalf("expected error for case %q", tt.name)
			}
		})
	}
}
