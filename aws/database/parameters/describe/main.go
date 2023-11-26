package describe

import (
	"fmt"

	"github.com/debeando/go-common/aws/rds"
	"github.com/debeando/go-common/log"
	"github.com/debeando/go-common/table"

	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "describe [IDENTIFIER]",
		Short: "Show all information about parameter group.",
		Example: `
  # Describe specific parameter group:
  zenit aws database parameters test-rds-parameter-group`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				cmd.Help()
				return
			}

			r := rds.RDS{}

			if err := r.Init(); err != nil {
				log.Error(err.Error())
				return
			}

			parameters, err := r.Parameters(args[0])
			if err != nil {
				fmt.Println(err)
				return
			}

			tbl := table.New()
			tbl.Column(0, table.Column{Name: "NAME"})
			tbl.Column(1, table.Column{Name: "VALUES"})
			tbl.Column(2, table.Column{Name: "APPLY METHOD"})
			tbl.Column(3, table.Column{Name: "APPLY TYPE"})
			tbl.Column(4, table.Column{Name: "MODIFIABLE"})
			for _, parameter := range parameters {
				tbl.Add(parameter.Name, parameter.Value, parameter.ApplyMethod, parameter.ApplyType, parameter.IsModifiable)
			}
			tbl.Print()
		},
	}

	return cmd
}
