package task

import (
	"fmt"

	"github.com/pragmaticbutton/todo/internal/domain/task"
	"github.com/pragmaticbutton/todo/internal/service"
	"github.com/spf13/cobra"
)

func NewAddCmd(taskService *service.TaskService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add <description>",
		Short: "Add a new task",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			description := args[0]
			_, err := taskService.AddTask(service.AddTaskInput{
				Description: description,
				Priority:    task.PriorityMedium,
				ListID:      nil,
			})
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Added task: %s\n", description)

			return nil
		},
	}
	return cmd
}
