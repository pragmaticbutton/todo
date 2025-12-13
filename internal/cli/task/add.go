package task

import (
	"fmt"

	"github.com/pragmaticbutton/todo/internal/domain/task"
	"github.com/pragmaticbutton/todo/internal/service"
	"github.com/spf13/cobra"
)

func NewAddCmd(taskService *service.TaskService) *cobra.Command {
	var (
		priorityStr string
		listIDRaw   uint32

		priority *task.Priority
		listID   *uint32
	)

	cmd := &cobra.Command{
		Use:   "add <description> [flags]",
		Short: "Add a new task",
		Args:  cobra.ExactArgs(1),

		RunE: func(cmd *cobra.Command, args []string) error {
			description := args[0]

			if cmd.Flags().Changed(flagPriority) {
				p, err := task.ParsePriority(priorityStr)
				if err != nil {
					return err
				}
				priority = &p
			}

			if cmd.Flags().Changed(flagListID) {
				id := listIDRaw
				listID = &id
			}

			_, err := taskService.AddTask(service.AddTaskInput{
				Description: description,
				Priority:    priority,
				ListID:      listID,
			})
			if err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Added task: %s\n", description)
			return nil
		},
	}

	// Optional flags (no CLI defaults — domain decides real defaults)
	cmd.Flags().StringVarP(
		&priorityStr,
		flagPriority,
		"p",
		"",
		"Optional priority: low, medium, high",
	)

	cmd.Flags().Uint32VarP(
		&listIDRaw,
		flagListID,
		"l",
		0,
		"Optional ID of the list to attach this task to",
	)

	return cmd
}
