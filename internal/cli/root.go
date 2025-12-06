package cli

import (
	"github.com/pragmaticbutton/todo/internal/cli/list"
	"github.com/pragmaticbutton/todo/internal/cli/task"
	"github.com/pragmaticbutton/todo/internal/service"
	"github.com/spf13/cobra"
)

func NewRootCmd(taskService *service.TaskService, listService *service.ListService) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "todo",
		Short:         "My TODO application",
		Long:          `The best TODO application in the world.`,
		SilenceUsage:  true,
		SilenceErrors: true,
		// Attach top-level groups here.
	}

	cmd.AddCommand(task.NewTaskCmd(taskService))
	cmd.AddCommand(list.NewListCmd(listService))

	return cmd
}
