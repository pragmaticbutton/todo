package task_test

import (
	"bytes"
	"path/filepath"
	"testing"

	taskcmd "github.com/pragmaticbutton/todo/internal/cli/task"
	"github.com/pragmaticbutton/todo/internal/domain/task"
	"github.com/pragmaticbutton/todo/internal/service"
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

func TestDeleteCmd_Integration(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		tasks   []task.Task
		args    []string
		checkFn func(t *testing.T, svc *service.TaskService)
	}{
		{
			name:  "delete single task",
			tasks: []task.Task{{ID: 1, Description: "A"}},
			args:  []string{"1"},
			checkFn: func(t *testing.T, svc *service.TaskService) {
				if _, err := svc.GetTask(1); err == nil {
					t.Fatalf("expected task to be deleted")
				}
			},
		},
		{
			name: "delete one task leaves others",
			tasks: []task.Task{
				{ID: 1, Description: "First"},
				{ID: 2, Description: "Second"},
			},
			args: []string{"2"},
			checkFn: func(t *testing.T, svc *service.TaskService) {
				if _, err := svc.GetTask(2); err == nil {
					t.Fatalf("expected deleted task to be missing")
				}
				remaining, err := svc.GetTask(1)
				if err != nil {
					t.Fatalf("expected other task to remain: %v", err)
				}
				if remaining.Description != "First" {
					t.Fatalf("unexpected remaining task: %+v", remaining)
				}
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := newTaskServiceWithTasks(t, tt.tasks...)
			cmd := taskcmd.NewDeleteCmd(svc)
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
		tasks []task.Task
		args  []string
	}{
		{
			name: "invalid id",
			args: []string{"abc"},
		},
		{
			name: "missing task",
			args: []string{"2"},
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

			svc := newTaskServiceWithTasks(t, tt.tasks...)
			cmd := taskcmd.NewDeleteCmd(svc)
			cmd.SetArgs(tt.args)

			if err := cmd.Execute(); err == nil {
				t.Fatalf("expected error for case %q", tt.name)
			}
		})
	}
}
