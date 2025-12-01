package task

import (
	"github.com/pragmaticbutton/todo/internal/service"
	"github.com/spf13/cobra"
)

func NewTaskCmd(taskService *service.TaskService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "task",
		Short: "Manage tasks",
	}

	cmd.AddCommand(
		NewAddCmd(taskService),
		NewListCmd(taskService),
		NewGetCmd(taskService),
		NewDeleteCmd(taskService),
		NewCompleteCmd(taskService),
		NewReopenCmd(taskService),
		NewUpdateCmd(taskService),
	)

	return cmd
}
