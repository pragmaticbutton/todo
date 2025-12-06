package list

import (
	"fmt"

	"github.com/pragmaticbutton/todo/internal/service"
	"github.com/spf13/cobra"
)

func NewAddCmd(listService *service.ListService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add <name> <description>",
		Short: "Add a new list",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			desc := args[1]

			if err := listService.AddList(service.AddListInput{
				Name:        name,
				Description: desc,
			}); err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Added list: %s\n", name)
			return nil
		},
	}

	return cmd
}
