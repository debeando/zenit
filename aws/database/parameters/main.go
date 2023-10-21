package parameters

import (
	"zenit/aws/database/parameters/describe"
	"zenit/aws/database/parameters/list"

	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "parameters",
		Short: "Manage parameter groups",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.AddCommand(list.NewCommand())
	cmd.AddCommand(describe.NewCommand())

	return cmd
}
