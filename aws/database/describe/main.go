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
		Short: "Show all information about database.",
		Example: `
  # Describe specific database instance:
  zenit aws database describe test-rds`,
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

			instance, err := r.Describe(args[0])
			if err != nil {
				fmt.Println(err)
				return
			}

			tbl := table.New()
			tbl.Column(0, table.Column{Name: "ATTRIBUTE", Alignment: table.Right})
			tbl.Column(1, table.Column{Name: "VALUE"})
			for k, v := range instance.JSON() {
				tbl.Add(k, v)
			}
			tbl.Print()
		},
	}

	return cmd
}
