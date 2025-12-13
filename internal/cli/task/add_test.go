package task_test

import (
	"bytes"
	"path/filepath"
	"testing"

	taskcmd "github.com/pragmaticbutton/todo/internal/cli/task"
	"github.com/pragmaticbutton/todo/internal/domain/task"
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
			args:   []string{"Bake cookies"},
			golden: filepath.Join("testdata", "add", "simple_add.golden"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := newTaskServiceWithTasks(t)
			cmd := taskcmd.NewAddCmd(svc)

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
		name        string
		description string
		args        []string
		checkFn     func(t *testing.T, tk *task.Task)
	}{
		{
			name:        "add task with description only",
			description: "Bake cookies",
			args:        []string{"Bake cookies"},
			checkFn: func(t *testing.T, tk *task.Task) {
				if tk.Description != "Bake cookies" {
					t.Fatalf("unexpected description: %s", tk.Description)
				}
				if tk.Priority != task.PriorityMedium {
					t.Fatalf("expected default PriorityMedium, got %v", tk.Priority)
				}
				if tk.ListID != nil {
					t.Fatalf("expected nil ListID, got %v", tk.ListID)
				}
			},
		},
		{
			name:        "add task with priority flag high",
			description: "Buy milk",
			args:        []string{"Buy milk", "-p", "high"},
			checkFn: func(t *testing.T, tk *task.Task) {
				if tk.Description != "Buy milk" {
					t.Fatalf("unexpected description: %s", tk.Description)
				}
				if tk.Priority != task.PriorityHigh {
					t.Fatalf("expected PriorityHigh, got %v", tk.Priority)
				}
			},
		},
		{
			name:        "add task with priority flag low",
			description: "Clean house",
			args:        []string{"Clean house", "-p", "low"},
			checkFn: func(t *testing.T, tk *task.Task) {
				if tk.Description != "Clean house" {
					t.Fatalf("unexpected description: %s", tk.Description)
				}
				if tk.Priority != task.PriorityLow {
					t.Fatalf("expected PriorityLow, got %v", tk.Priority)
				}
			},
		},
		{
			name:        "add task with priority flag medium",
			description: "Medium task",
			args:        []string{"Medium task", "-p", "medium"},
			checkFn: func(t *testing.T, tk *task.Task) {
				if tk.Description != "Medium task" {
					t.Fatalf("unexpected description: %s", tk.Description)
				}
				if tk.Priority != task.PriorityMedium {
					t.Fatalf("expected PriorityMedium, got %v", tk.Priority)
				}
			},
		},
		{
			name:        "add task with long priority flag",
			description: "Long flag task",
			args:        []string{"Long flag task", "--priority", "high"},
			checkFn: func(t *testing.T, tk *task.Task) {
				if tk.Priority != task.PriorityHigh {
					t.Fatalf("expected PriorityHigh, got %v", tk.Priority)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := newTaskServiceWithTasks(t)
			cmd := taskcmd.NewAddCmd(svc)
			cmd.SetArgs(tt.args)

			if err := cmd.Execute(); err != nil {
				t.Fatalf("execute add: %v", err)
			}

			tasks, err := svc.ListTasks()
			if err != nil {
				t.Fatalf("list tasks: %v", err)
			}
			if len(tasks) != 1 {
				t.Fatalf("expected 1 task, got %d", len(tasks))
			}

			tt.checkFn(t, &tasks[0])
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
			name: "invalid priority value",
			args: []string{"Task", "-p", "invalid"},
		},
		{
			name: "missing description",
			args: []string{},
		},
		{
			name: "invalid list-id value",
			args: []string{"Task", "-l", "abc"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := newTaskServiceWithTasks(t)
			cmd := taskcmd.NewAddCmd(svc)
			cmd.SetArgs(tt.args)

			if err := cmd.Execute(); err == nil {
				t.Fatalf("expected error for case %q", tt.name)
			}
		})
	}
}
