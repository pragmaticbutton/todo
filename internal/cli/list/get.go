package list

import (
	"fmt"
	"strconv"

	"github.com/pragmaticbutton/todo/internal/service"
	"github.com/spf13/cobra"
)

func NewGetCmd(listService *service.ListService) *cobra.Command {
	return &cobra.Command{
		Use:   "get <list-id>",
		Short: "Get a specific list",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid list id")
			}
			l, err := listService.GetList(uint32(id))
			if err != nil {
				return err
			}

			out := cmd.OutOrStdout()
			fmt.Fprintln(out, "ID   Name        Description")
			fmt.Fprintln(out, "--   ----        -----------")
			fmt.Fprintf(out, "%-4d %-12s %s\n", l.ID, l.Name, l.Description)

			return nil
		},
	}
}
