package databases

import (
	"github.com/debeando/go-common/aws/rds"
	"github.com/debeando/go-common/log"
	"github.com/debeando/go-common/table"

	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "databases",
		Short: "List all databases.",
		Run: func(cmd *cobra.Command, args []string) {
			r := rds.RDS{}

			if err := r.Init(); err != nil {
				log.Error(err.Error())
				return
			}

			instances := r.List()

			tbl := table.New()
			tbl.Column(0, table.Column{Name: "ENGINE"})
			tbl.Column(1, table.Column{Name: "VERSION"})
			tbl.Column(2, table.Column{Name: "IDENTIFIER"})
			tbl.Column(3, table.Column{Name: "CLASS"})
			tbl.Column(4, table.Column{Name: "STATUS"})
			for _, instance := range instances {
				tbl.Add(instance.Engine, instance.Version, instance.Identifier, instance.Class, instance.Status)
			}
			tbl.Print()
		},
	}

	return cmd
}
