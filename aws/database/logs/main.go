package logs

import (
	"zenit/aws/database/logs/list"
	"zenit/aws/database/logs/view"

	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "logs",
		Short: "Logs about database; error, slow logs, general log, ...",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.AddCommand(list.NewCommand())
	cmd.AddCommand(view.NewCommand())

	return cmd
}
