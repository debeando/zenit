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

			tbl := table.New("NAME", "VALUES", "APPLY METHOD", "APPLY TYPE", "MODIFIABLE")
			for _, parameter := range parameters {
				tbl.AddRow(parameter.Name, parameter.Value, parameter.ApplyMethod, parameter.ApplyType, parameter.IsModifiable)
			}
			tbl.Print()
		},
	}

	return cmd
}
