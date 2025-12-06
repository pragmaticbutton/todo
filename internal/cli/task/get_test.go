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
