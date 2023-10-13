package describe

import (
	"fmt"
	"encoding/json"

	"github.com/debeando/go-common/aws/rds"

	"github.com/rodaine/table"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "describe [IDENTIFIER]",
		Short: "Show all information about database.",
		Example: `
  # Describe specific database instance
  zenit aws database describe test-rds`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 1 {
				cmd.Help()
				return
			}

			r := rds.Config{}
			r.Init()
			i, err := r.Describe(args[0])
			if err != nil{
				fmt.Println(err)
				return
			}

			var z map[string]interface{}
			a, _ := json.Marshal(i)
			json.Unmarshal(a, &z)

			tbl := table.New("ATTRIBUTE", "VALUE")
			for k,v := range z {
				tbl.AddRow(k,v)
			}
			tbl.Print()
		},
	}

	return cmd
}
