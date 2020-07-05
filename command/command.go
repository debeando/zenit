package command

import (
	"flag"
	"fmt"
	"os"

	"github.com/swapbyt3s/zenit/command/daemon"
	"github.com/swapbyt3s/zenit/command/example"
	"github.com/swapbyt3s/zenit/common/log"
	"github.com/swapbyt3s/zenit/config"
)

// USAGE is a const to have help description for CLI.
const USAGE = `zenit %s, agent for collecting and reporting metrics for
MySQL, Percona, MariaDB and ProxySQL.

Usage: %s [--help | --install | --uninstall | --version]

Options:
  --help        Show this help.
  --config      Print out full sample configuration to stdout.
  --debug       Enable debug mode.
  --install     Install service on system.
  --uninstall   Uninstall service on system.
  --version     Print version numbers.

Example:

  # Generate a zenit config file:
  $ sudo zenit --config > /etc/zenit/zenit.yaml

  # Install init script:
  $ sudo zenit --install

  # Start zenit agent:
  $ sudo [systemctl|initctl] [start|stop|restart|status] zenit

For more help, plese visit: https://github.com/swapbyt3s/zenit/wiki

`

func init() {
	if err := config.File.Load(); err != nil {
		log.Error("Config", map[string]interface{}{"error": err})
		os.Exit(1)
	}
	
	if warn := config.File.SanityCheck(); len(warn) > 0 {
		log.Warning("Config", map[string]interface{}{"message": warn})
	}
}

func Run() {
	fHelp      := flag.Bool("help", false, "Show this help.")
	fConfig    := flag.Bool("config", false, "Print out full sample configuration to stdout.")
	fInstall   := flag.Bool("install", false, "Install service on system.")
	fUninstall := flag.Bool("uninstall", false, "Uninstall service on system.")
	fVersion   := flag.Bool("version", false, "Show version.")
	_           = flag.Bool("debug", false, "Enable debug mode.")

	flag.Usage = func() { help(1) }
	flag.Parse()

	log.Configure()
	daemon.Configure()

	switch {
	case *fVersion:
		fmt.Println(Version())
	case *fHelp:
		help(0)
	case *fConfig:
		fmt.Printf(example.GetConfigFile())
	case *fInstall:
		daemon.Install()
	case *fUninstall:
		daemon.Uninstall()
	default:
		daemon.Daemonize()
	}
}

func help(rc int) {
	fmt.Printf(USAGE, Version(), os.Args[0])
	os.Exit(rc)
}
