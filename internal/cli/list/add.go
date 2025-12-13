package list

import (
	"fmt"

	"github.com/pragmaticbutton/todo/internal/service"
	"github.com/spf13/cobra"
)

const flagDescription = "description"

func NewAddCmd(listService *service.ListService) *cobra.Command {
	var description string

	cmd := &cobra.Command{
		Use:   "add <name> [flags]",
		Short: "Add a new list",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			desc := ""
			if cmd.Flags().Changed(flagDescription) {
				desc = description
			}

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

	cmd.Flags().StringVarP(
		&description,
		flagDescription,
		"d",
		"",
		"Optional list description",
	)

	return cmd
}
