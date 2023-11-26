package list

import (
	"fmt"

	"github.com/debeando/go-common/aws/rds"
	"github.com/debeando/go-common/log"
	"github.com/debeando/go-common/table"

	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "list",
		Short: "List parameter group available.",
		Example: `
  # List parameter groups available:
  zenit aws database parameters list`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				cmd.Help()
				return
			}

			r := rds.RDS{}

			if err := r.Init(); err != nil {
				log.Error(err.Error())
				return
			}

			parameters, err := r.ParametersGroup()
			if err != nil {
				fmt.Println(err)
				return
			}

			tbl := table.New()
			tbl.Column(0, table.Column{Name: "NAME"})
			tbl.Column(1, table.Column{Name: "DESCRIPTION"})
			tbl.Column(2, table.Column{Name: "FAMILY"})
			for _, parameter := range parameters {
				tbl.Add(parameter.Name, parameter.Description, parameter.Family)
			}
			tbl.Print()

			fmt.Println()

		},
	}

	return cmd
}
