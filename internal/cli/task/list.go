package task

import (
	"fmt"

	"github.com/pragmaticbutton/todo/internal/service"
	"github.com/spf13/cobra"
)

func NewListCmd(taskService *service.TaskService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all tasks",
		RunE: func(cmd *cobra.Command, args []string) error {
			ts, err := taskService.ListTasks()
			if err != nil {
				return err
			}
			out := cmd.OutOrStdout()

			fmt.Fprintln(out, "    ID  Description")
			fmt.Fprintln(out, "    --  -----------")
			for _, t := range ts {
				status := " "
				if t.Done {
					status = "x"
				}
				fmt.Fprintf(out, "[%s] %-4d %s\n", status, t.ID, t.Description)
			}
			return nil
		},
	}
	return cmd
}
