package agent

import (
	"fmt"
	"os"

	"zenit/config"
	"zenit/config/example"
	"zenit/agent/plugins"

	"github.com/debeando/go-common/file"
	"github.com/debeando/go-common/log"
	"github.com/debeando/go-common/net"
	"github.com/spf13/cobra"
)

var configPath string
var configExample bool

func NewCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "agent",
		Short: "Agent to collect any information for many services; MySQL, ProxySQL and more...",
		Example: `
  Generate a zenit config file:
    $ zenit agent --config-example > /etc/zenit/zenit.yaml

  Use specific config file:
    $ zenit agent --config=/etc/zenit/zenit.yaml`,
		Run: func(cmd *cobra.Command, args []string) {
			if configExample {
				fmt.Print(example.Load())
				os.Exit(0)
			}

			if file.Exist(configPath) {
				c := config.GetInstance()
				c.Path = configPath
				c.IPAddress = net.IPAddress()

				c.Load()
				plugins.Load()
			} else {
				log.Error("Invalid config file, please verify.")
			}
		},
	}

	cmd.Flags().StringVar(&configPath, "config", "zenit.yaml", "Config path")
	cmd.Flags().BoolVar(&configExample, "config-example", false, "Print out full sample configuration to stdout in YAML format.")

	return cmd
}
