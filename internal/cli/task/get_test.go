package task_test

import (
	"bytes"
	"path/filepath"
	"strconv"
	"testing"

	taskcmd "github.com/pragmaticbutton/todo/internal/cli/task"
	"github.com/pragmaticbutton/todo/internal/domain/task"
)

func TestGetCmd_Golden(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		task   task.Task
		golden string
	}{
		{
			name: "pending task",
			task: task.Task{
				ID:          1,
				Description: "Bake cookies",
				Done:        false,
			},
			golden: "pending.golden",
		},
		{
			name: "completed task",
			task: task.Task{
				ID:          2,
				Description: "Ship the tests",
				Done:        true,
			},
			golden: "completed.golden",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			taskService := newTaskServiceWithTasks(t, tc.task)
			cmd := taskcmd.NewGetCmd(taskService)

			var buf bytes.Buffer
			cmd.SetOut(&buf)
			cmd.SetArgs([]string{strconv.Itoa(int(tc.task.ID))})

			if err := cmd.Execute(); err != nil {
				t.Fatalf("execute get command: %v", err)
			}

			goldenPath := filepath.Join("testdata", "get", tc.golden)
			assertGolden(t, buf.String(), goldenPath)
		})
	}
}

func TestGetCmd_Integration(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		tasks   []task.Task
		args    []string
		checkFn func(t *testing.T, tk *task.Task)
	}{
		{
			name:  "gets existing task",
			tasks: []task.Task{{ID: 1, Description: "Bake cookies"}},
			args:  []string{"1"},
			checkFn: func(t *testing.T, tk *task.Task) {
				if tk.Description != "Bake cookies" || tk.Done {
					t.Fatalf("unexpected task %+v", tk)
				}
			},
		},
		{
			name: "gets correct task when multiple exist",
			tasks: []task.Task{
				{ID: 1, Description: "First"},
				{ID: 2, Description: "Second", Done: true},
			},
			args: []string{"2"},
			checkFn: func(t *testing.T, tk *task.Task) {
				if tk.Description != "Second" || !tk.Done {
					t.Fatalf("expected completed second task, got %+v", tk)
				}
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := newTaskServiceWithTasks(t, tt.tasks...)
			cmd := taskcmd.NewGetCmd(svc)
			cmd.SetArgs(tt.args)

			if err := cmd.Execute(); err != nil {
				t.Fatalf("execute get command: %v", err)
			}

			result, err := svc.GetTask(parseID(t, tt.args[0]))
			if err != nil {
				t.Fatalf("get task: %v", err)
			}

			tt.checkFn(t, result)
		})
	}
}

func TestGetCmd_Errors(t *testing.T) {
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
			cmd := taskcmd.NewGetCmd(svc)
			cmd.SetArgs(tt.args)

			if err := cmd.Execute(); err == nil {
				t.Fatalf("expected error for case %q", tt.name)
			}
		})
	}
}

func parseID(t *testing.T, raw string) uint32 {
	t.Helper()

	id, err := strconv.Atoi(raw)
	if err != nil {
		t.Fatalf("parse id: %v", err)
	}
	return uint32(id)
}
