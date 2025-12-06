package list

import (
	"fmt"
	"strconv"

	"github.com/pragmaticbutton/todo/internal/service"
	"github.com/spf13/cobra"
)

func NewUpdateCmd(listService *service.ListService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <list-id> <name> <description>",
		Short: "Update a list",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid list id")
			}

			name := args[1]
			desc := args[2]
			input := service.UpdateListInput{
				Name:        &name,
				Description: &desc,
			}

			if err := listService.UpdateList(uint32(id), &input); err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Updated list %d\n", id)
			return nil
		},
	}

	return cmd
}
