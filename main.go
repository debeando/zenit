package main

import (
	"zenit/agent"
	"zenit/aws"
	"zenit/mysql"
	"zenit/service"
	"zenit/version"

	"github.com/debeando/go-common/log"

	"github.com/spf13/cobra"
)

var debug bool

func main() {
	var rootCmd = &cobra.Command{
		Use: "zenit [COMMANDS] [OPTIONS]",
		Long: `zenit is a multipurpose tool for a MySQL, you can; monitoring, lint
data model and more, please see available commands.

Find more information at: https://github.com/debeando/zenit`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if debug {
				log.SetLevel(log.DebugLevel)
			}
		},
	}

	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "v", false, "Debug output")

	rootCmd.AddCommand(agent.NewCommand())
	rootCmd.AddCommand(aws.NewCommand())
	rootCmd.AddCommand(mysql.NewCommand())
	rootCmd.AddCommand(service.NewCommand())
	rootCmd.AddCommand(version.NewCommand())
	rootCmd.Execute()
}
