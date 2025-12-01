package task

import (
	"fmt"
	"strconv"

	"github.com/pragmaticbutton/todo/internal/service"
	"github.com/spf13/cobra"
)

func NewUpdateCmd(taskService *service.TaskService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <task-id> <description>",
		Short: "Update a task description",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid task id")
			}
			desc := args[1]
			input := service.UpdateTaskInput{Description: &desc}

			if _, err := taskService.UpdateTask(uint32(id), input); err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Updated task %d\n", id)
			return nil
		},
	}
	return cmd
}
