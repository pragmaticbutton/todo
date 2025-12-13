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
			args:   []string{"1", "--description", "NewDesc"},
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
	cmd.SetArgs([]string{"1", "--description", "NewDesc"})

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

func TestUpdateCmd_Integration(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		task    task.Task
		args    []string
		checkFn func(t *testing.T, tk *task.Task)
	}{
		{
			name: "update priority",
			task: task.Task{ID: 1, Description: "Old", Priority: task.PriorityLow},
			args: []string{"1", "--priority", "high"},
			checkFn: func(t *testing.T, tk *task.Task) {
				if tk.Priority != task.PriorityHigh {
					t.Fatalf("expected priority to be high, got %v", tk.Priority)
				}
			},
		},
		{
			name: "update done true",
			task: task.Task{ID: 1, Description: "Old"},
			args: []string{"1", "--done"},
			checkFn: func(t *testing.T, tk *task.Task) {
				if !tk.Done {
					t.Fatalf("expected task to be marked done")
				}
			},
		},
		{
			name: "update done false",
			task: task.Task{ID: 1, Description: "Old", Done: true},
			args: []string{"1", "--done=false"},
			checkFn: func(t *testing.T, tk *task.Task) {
				if tk.Done {
					t.Fatalf("expected task to be reopened")
				}
			},
		},
		{
			name: "update list id",
			task: task.Task{ID: 1, Description: "Old"},
			args: []string{"1", "--list-id", "3"},
			checkFn: func(t *testing.T, tk *task.Task) {
				if tk.ListID == nil || *tk.ListID != 3 {
					t.Fatalf("expected list id 3, got %v", tk.ListID)
				}
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := newTaskServiceWithTasks(t, tt.task)
			cmd := taskcmd.NewUpdateCmd(svc)
			cmd.SetArgs(tt.args)

			if err := cmd.Execute(); err != nil {
				t.Fatalf("execute update: %v", err)
			}

			updated, err := svc.GetTask(1)
			if err != nil {
				t.Fatalf("get task: %v", err)
			}

			tt.checkFn(t, updated)
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
			args: []string{"abc", "--description", "NewDesc"},
		},
		{
			name: "missing task",
			args: []string{"2", "--description", "NewDesc"},
		},
		{
			name: "invalid priority",
			args: []string{"1", "--priority", "urgent"},
		},
		{
			name: "no fields provided",
			args: []string{"1"},
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
