package command

import (
	"flag"
	"fmt"
	"os"
)

// USAGE is a const to have help description for CLI.
const USAGE = `zenit %s
Usage: %s [--help | --install | --uninstall | --version]
Options:
  --help        Show this help.
  --install     Install service on system.
  --uninstall   Uninstall service on system.
  --version     Print version numbers.
`

func Run() {
	fHelp      := flag.Bool("help", false, "Show this help.")
	fInstall   := flag.Bool("install", false, "Install service on system.")
	fUninstall := flag.Bool("uninstall", false, "Uninstall service on system.")
	fVersion   := flag.Bool("version", false, "Show version.")

	flag.Usage = func() { help(1) }
	flag.Parse()

	switch {
	case *fVersion:
		fmt.Println(Version())
	case *fHelp:
		help(0)
	case *fInstall:
		Install()
	case *fUninstall:
		Uninstall()
	default:
		Daemonize()
	}
}

func help(rc int) {
	fmt.Printf(USAGE, Version(), os.Args[0])
	os.Exit(rc)
}
