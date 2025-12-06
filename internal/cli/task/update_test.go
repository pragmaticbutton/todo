package task_test

import (
	"bytes"
	"path/filepath"
	"testing"

	taskcmd "github.com/pragmaticbutton/todo/internal/cli/task"
	"github.com/pragmaticbutton/todo/internal/domain/task"
)

func TestUpdateCmd_Golden(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		task   task.Task
		args   []string
		golden string
	}{
		{
			name:   "update description",
			task:   task.Task{ID: 1, Description: "Old"},
			args:   []string{"1", "NewDesc"},
			golden: filepath.Join("testdata", "update", "updated.golden"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := newTaskServiceWithTasks(t, tt.task)
			cmd := taskcmd.NewUpdateCmd(svc)

			var buf bytes.Buffer
			cmd.SetOut(&buf)
			cmd.SetArgs(tt.args)

			if err := cmd.Execute(); err != nil {
				t.Fatalf("execute update: %v", err)
			}

			assertGolden(t, buf.String(), tt.golden)
		})
	}
}

func TestUpdateCmd_Unit(t *testing.T) {
	t.Parallel()

	svc := newTaskServiceWithTasks(t, task.Task{ID: 1, Description: "Old"})
	cmd := taskcmd.NewUpdateCmd(svc)
	cmd.SetArgs([]string{"1", "NewDesc"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute update: %v", err)
	}

	updated, err := svc.GetTask(1)
	if err != nil {
		t.Fatalf("get task: %v", err)
	}
	if updated.Description != "NewDesc" {
		t.Fatalf("expected description to be updated, got %s", updated.Description)
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
			args: []string{"abc", "NewDesc"},
		},
		{
			name: "missing task",
			args: []string{"2", "NewDesc"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := newTaskServiceWithTasks(t, task.Task{ID: 1, Description: "Old"})
			cmd := taskcmd.NewUpdateCmd(svc)
			cmd.SetArgs(tt.args)

			if err := cmd.Execute(); err == nil {
				t.Fatalf("expected error for case %q", tt.name)
			}
		})
	}
}
