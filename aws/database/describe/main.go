package describe

import (
	"fmt"

	"github.com/debeando/go-common/aws/rds"
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

			r := rds.Config{}
			r.Init()
			instance, err := r.Describe(args[0])
			if err != nil {
				fmt.Println(err)
				return
			}

			tbl := table.New("ATTRIBUTE", "VALUE")
			tbl.SetFirstColumnAlignment(table.Right)
			for k, v := range instance.JSON() {
				tbl.AddRow(k, v)
			}
			tbl.Print()
		},
	}

	return cmd
}