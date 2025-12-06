package task_test

import (
	"bytes"
	"path/filepath"
	"testing"

	taskcmd "github.com/pragmaticbutton/todo/internal/cli/task"
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

func TestAddCmd_Unit(t *testing.T) {
	t.Parallel()

	svc := newTaskServiceWithTasks(t)
	cmd := taskcmd.NewAddCmd(svc)
	cmd.SetArgs([]string{"Bake cookies"})

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
	if tasks[0].Description != "Bake cookies" {
		t.Fatalf("unexpected task description: %s", tasks[0].Description)
	}
}
