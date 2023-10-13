package databases

import (
	"github.com/debeando/go-common/aws/rds"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "databases",
		Short: "List all databases.",
		Run: func(cmd *cobra.Command, args []string) {
			r := rds.Config{}
			r.Init()
			instances := r.List()

			tbl := table.New("ENGINE", "VERSION", "IDENTIFIER", "CLASS", "STATUS")
			for _, instance := range instances {
				tbl.AddRow(instance.Engine, instance.Version, instance.Identifier, instance.Class, instance.Status)
			}
			tbl.Print()
		},
	}

	return cmd
}
