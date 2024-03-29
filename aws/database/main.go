package database

import (
	"zenit/aws/database/describe"
	"zenit/aws/database/logs"
	"zenit/aws/database/parameters"

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
	cmd.AddCommand(logs.NewCommand())
	cmd.AddCommand(parameters.NewCommand())

	return cmd
}
