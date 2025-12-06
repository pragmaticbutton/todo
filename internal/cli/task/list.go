package task

import (
	"fmt"
	"sort"

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
			sort.Slice(ts, func(i, j int) bool { return ts[i].ID < ts[j].ID })

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
