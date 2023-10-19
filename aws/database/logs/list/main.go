package list

import (
	"fmt"

	"github.com/debeando/go-common/aws/rds"
	"github.com/debeando/go-common/log"
	"github.com/debeando/go-common/table"

	"github.com/spf13/cobra"
)

var filter string
var sort bool

func NewCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "list [IDENTIFIER] [FILENAME]",
		Short: "List all logs about database.",
		Example: `
  # List logs of specific database instance:
  zenit aws database logs test-rds

  # List slow query logs of specific database instance:
  zenit aws database logs test-rds --filter=slowquery

  # Sort list slow query logs of specific database instance:
  zenit aws database logs test-rds --filter=slowquery --sort`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 || len(args) > 1 {
				cmd.Help()
				return
			}

			r := rds.RDS{}

			if err := r.Init(); err != nil {
				log.Error(err.Error())
				return
			}

			logs, err := r.Logs(args[0], filter)
			if err != nil {
				fmt.Println(err)
				return
			}

			if sort {
				logs.SortBySize()
			}

			tbl := table.New("FILE", "SIZE")
			for _, log := range logs {
				tbl.AddRow(log.FileName, log.Size)
			}
			tbl.Print()
		},
	}

	cmd.Flags().StringVar(&filter, "filter", "", "Filters the available log files for log file names that contain the specified string.")
	cmd.Flags().BoolVar(&sort, "sort", false, "Sort list of logs by file size.")

	return cmd
}
