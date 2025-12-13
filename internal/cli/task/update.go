package task

import (
	"fmt"
	"strconv"

	"github.com/pragmaticbutton/todo/internal/domain/task"
	"github.com/pragmaticbutton/todo/internal/service"
	"github.com/spf13/cobra"
)

func NewUpdateCmd(taskService *service.TaskService) *cobra.Command {
	var (
		description string
		priorityStr string
		doneRaw     bool
		listIDRaw   uint32
	)

	cmd := &cobra.Command{
		Use:   "update <task-id> [flags]",
		Short: "Update a task",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid task id")
			}

			input := service.UpdateTaskInput{}

			if cmd.Flags().Changed(flagDescription) {
				desc := description
				input.Description = &desc
			}
			if cmd.Flags().Changed(flagPriority) {
				p, err := task.ParsePriority(priorityStr)
				if err != nil {
					return err
				}
				input.Priority = &p
			}
			if cmd.Flags().Changed(flagDone) {
				done := doneRaw
				input.Done = &done
			}
			if cmd.Flags().Changed(flagListID) {
				lID := listIDRaw
				input.ListID = &lID
			}

			if input.Description == nil && input.Priority == nil && input.Done == nil && input.ListID == nil {
				return fmt.Errorf("no fields provided to update")
			}

			if _, err := taskService.UpdateTask(uint32(id), input); err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Updated task %d\n", id)
			return nil
		},
	}

	cmd.Flags().StringVarP(&description, flagDescription, "d", "", "Optional new description")
	cmd.Flags().StringVarP(&priorityStr, flagPriority, "p", "", "Optional priority: low, medium, high")
	cmd.Flags().BoolVar(&doneRaw, flagDone, false, "Optional completion status")
	cmd.Flags().Uint32VarP(&listIDRaw, flagListID, "l", 0, "Optional list id to attach this task to")
	return cmd
}
