package cli

import (
	"github.com/pragmaticbutton/todo/internal/cli/interactive"
	"github.com/pragmaticbutton/todo/internal/cli/list"
	"github.com/pragmaticbutton/todo/internal/cli/task"
	"github.com/pragmaticbutton/todo/internal/service"
	"github.com/spf13/cobra"
)

func NewRootCmd(taskService *service.TaskService, listService *service.ListService) *cobra.Command {

	var interactiveMode bool

	cmd := &cobra.Command{
		Use:   "todo",
		Short: "My TODO application",
		Long:  `The best TODO application in the world.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if interactiveMode {
				return interactive.Run(cmd)
			}

			return cmd.Help()
		},
	}

	cmd.Flags().BoolVarP(
		&interactiveMode,
		"interactive",
		"i",
		false,
		"Run in interactive mode",
	)

	cmd.AddCommand(task.NewTaskCmd(taskService))
	cmd.AddCommand(list.NewListCmd(listService))

	return cmd
}
