package task

import (
	"fmt"
	"strconv"

	"github.com/pragmaticbutton/todo/internal/service"
	"github.com/spf13/cobra"
)

func NewReopenCmd(taskService *service.TaskService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reopen <task-id>",
		Short: "Reopen a completed task",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid task id")
			}
			if err := taskService.ReopenTask(uint32(id)); err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Reopened task %d\n", id)
			return nil
		},
	}
	return cmd
}
