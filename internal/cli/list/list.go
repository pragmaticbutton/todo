package list

import (
	"fmt"
	"sort"

	"github.com/pragmaticbutton/todo/internal/service"
	"github.com/spf13/cobra"
)

func NewListCmd(listService *service.ListService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Manage lists",
	}

	cmd.AddCommand(
		NewAddCmd(listService),
		NewListListsCmd(listService),
		NewGetCmd(listService),
		NewDeleteCmd(listService),
		NewUpdateCmd(listService),
		NewListTasksCmd(listService),
	)

	return cmd
}

func NewListListsCmd(listService *service.ListService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all lists",
		RunE: func(cmd *cobra.Command, args []string) error {
			ls, err := listService.ListLists()
			if err != nil {
				return err
			}
			sort.Slice(ls, func(i, j int) bool { return ls[i].ID < ls[j].ID })
			out := cmd.OutOrStdout()

			fmt.Fprintln(out, "ID   Name        Description")
			fmt.Fprintln(out, "--   ----        -----------")
			for _, l := range ls {
				fmt.Fprintf(out, "%-4d %-12s %s\n", l.ID, l.Name, l.Description)
			}
			return nil
		},
	}

	return cmd
}
