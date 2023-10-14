package logs

import (
	"fmt"

	"github.com/debeando/go-common/aws/rds"
	"github.com/debeando/go-common/table"

	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "logs [IDENTIFIER] [FILENAME]",
		Short: "List all logs about database.",
		Example: `
  # List logs of specific database instance:
  zenit aws database logs test-rds

  # Show logs to specific path and file name:
  zenit aws database logs test-rds error/mysql-error-running.log`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 || len(args) > 2 {
				cmd.Help()
				return
			}

			r := rds.Config{}
			r.Init()

			if len(args) == 1 && len(args[0]) > 0 {
				logs, err := r.Logs(args[0])
				if err != nil {
					fmt.Println(err)
					return
				}
				tbl := table.New("FILE", "SIZE")
				for _, log := range logs {
					tbl.AddRow(log.FileName, log.Size)
				}
				tbl.Print()
			}

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
