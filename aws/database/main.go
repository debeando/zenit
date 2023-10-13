package database

import (
	"zenit/aws/database/describe"

	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "database",
		Short: "Commands for specific database.",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.AddCommand(describe.NewCommand())

	return cmd
}
