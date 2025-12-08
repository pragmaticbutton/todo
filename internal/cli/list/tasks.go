package list

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/pragmaticbutton/todo/internal/service"
	"github.com/spf13/cobra"
)

func NewListTasksCmd(listService *service.ListService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tasks <list-id>",
		Short: "List tasks for a specific list",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid list id")
			}

			ts, err := listService.ListTasks(uint32(id))
			if err != nil {
				return err
			}
			sort.Slice(ts, func(i, j int) bool { return ts[i].ID < ts[j].ID })

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
