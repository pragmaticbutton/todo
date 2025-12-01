package task

import (
	"fmt"
	"strconv"

	"github.com/pragmaticbutton/todo/internal/service"
	"github.com/spf13/cobra"
)

func NewGetCmd(taskService *service.TaskService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <task-id>",
		Short: "Get a task",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid task id")
			}
			t, err := taskService.GetTask(uint32(id))
			if err != nil {
				return err
			}

			out := cmd.OutOrStdout()
			status := " "
			if t.Done {
				status = "x"
			}
			fmt.Fprintln(out, "    ID  Description")
			fmt.Fprintln(out, "    --  -----------")
			fmt.Fprintf(out, "[%s] %-4d %s\n", status, t.ID, t.Description)

			return nil
		},
	}
	return cmd
}
