package aws

import (
	"zenit/aws/database"
	"zenit/aws/databases"

	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "aws",
		Short: "AWS commands to facilitate administration.",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.AddCommand(database.NewCommand())
	cmd.AddCommand(databases.NewCommand())

	return cmd
}
