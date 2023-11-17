package mysql

import (
	"zenit/mysql/digest"
	"zenit/mysql/top"

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
	cmd.AddCommand(top.NewCommand())

	return cmd
}
