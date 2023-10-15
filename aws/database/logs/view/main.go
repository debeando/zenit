package view

import (
	"fmt"

	"github.com/debeando/go-common/aws/rds"

	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "view [IDENTIFIER] [FILENAME]",
		Short: "View logs details of specific database instance.",
		Example: `
  # View logs details to specific path and file name:
  zenit aws database logs view test-rds error/mysql-error-running.log`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 2 {
				cmd.Help()
				return
			}

			r := rds.Config{}
			r.Init()

			if len(args) == 2 && len(args[0]) > 0 && len(args[1]) > 0 {
				data, err := r.PollLogs(args[0], args[1])
				if err != nil {
					fmt.Println(err)
					return
				}

				if len(data) > 0 {
					fmt.Println(data)
				}
			}
		},
	}

	return cmd
}
