package list

import (
	"fmt"
	"strconv"

	"github.com/pragmaticbutton/todo/internal/service"
	"github.com/spf13/cobra"
)

func NewUpdateCmd(listService *service.ListService) *cobra.Command {
	var (
		name        string
		description string
	)
	cmd := &cobra.Command{
		Use:   "update <list-id> [flags]",
		Short: "Update a list",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid list id")
			}

			input := service.UpdateListInput{}

			if cmd.Flags().Changed(flagName) {
				n := name
				input.Name = &n
			}
			if cmd.Flags().Changed(flagDescription) {
				d := description
				input.Description = &d
			}

			// rethink when working on field validation
			if input.Name == nil && input.Description == nil {
				return fmt.Errorf("no fields provided to update")
			}

			if err := listService.UpdateList(uint32(id), &input); err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Updated list %d\n", id)
			return nil
		},
	}

	cmd.Flags().StringVarP(&name, flagName, "n", "", "Optional new name")
	cmd.Flags().StringVarP(&description, flagDescription, "d", "", "Optional new description")

	return cmd
}
