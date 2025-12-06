package list

import (
	"fmt"
	"strconv"

	"github.com/pragmaticbutton/todo/internal/service"
	"github.com/spf13/cobra"
)

func NewDeleteCmd(listService *service.ListService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <list-id>",
		Short: "Delete a list",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid list id")
			}

			if err := listService.DeleteList(uint32(id)); err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Deleted list %d\n", id)
			return nil
		},
	}

	return cmd
}
