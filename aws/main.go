package aws

import (
	"github.com/debeando/go-common/aws/rds"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	var cmdDatabases = &cobra.Command{
		Use:   "databases",
		Short: "List all databases.",
		Run: func(cmd *cobra.Command, args []string) {
			r := rds.Config{}
			r.Init()
			instances := r.List()

			tbl := table.New("Engine", "Version", "Identifier", "Class", "Status")
			for _, instance := range instances {
				tbl.AddRow(instance.Engine, instance.Version, instance.Identifier, instance.Class, instance.Status)
			}
			tbl.Print()
		},
	}

	var cmd = &cobra.Command{
		Use:   "aws",
		Short: "AWS commands to facilitate administration.",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	cmd.AddCommand(cmdDatabases)

	return cmd
}
