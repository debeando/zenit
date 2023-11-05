package mysql

import (
	"zenit/mysql/digest"

	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "mysql",
		Short: "MySQL commands to facilitate administration.",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.AddCommand(digest.NewCommand())

	return cmd
}
