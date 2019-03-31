package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins"

	"github.com/kardianos/service"
)

// USAGE is a const to have help description for CLI.
const USAGE = `zenit (%s) written by %s
Usage: %s [--help | --install | --uninstall | --version]
Options:
  --help        Show this help.
  --install     Install service on system.
  --uninstall   Uninstall service on system.
  --version     Print version numbers.
`

var (
	logger service.Logger

	// Build date and time when building.
	build  string
)

type program struct{}

func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	go p.run()
	return nil
}

func (p *program) run() {
	var wg sync.WaitGroup

	wg.Add(1)

	go plugins.Load(&wg)

	wg.Wait()
}

func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	return nil
}

func init() {
	config.Load()
	config.SanityCheck()
}

func main() {
	svcConfig := &service.Config{
		Name:        "zenit",
		DisplayName: "Zenit",
		Description: "Zenit Agent",
		Executable:  "/usr/bin/zenit",
	}

	prg := &program{}

	fHelp := flag.Bool("help", false, "Show this help.")
	fInstall := flag.Bool("install", false, "Install service on system.")
	fUninstall := flag.Bool("uninstall", false, "Uninstall service on system.")
	fVersion := flag.Bool("version", false, "Show version.")

	flag.Usage = func() { help(1) }
	flag.Parse()

	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	logger, err = s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}

	switch {
	case *fVersion:
		fmt.Printf("%s (%s)\n", config.Version, build)
		os.Exit(0)
	case *fHelp:
		help(0)
	case *fInstall:
		s.Install()
		os.Exit(0)
	case *fUninstall:
		s.Uninstall()
		os.Exit(0)
	}

	err = s.Run()
	if err != nil {
		logger.Error(err)
	}
}

func help(rc int) {
	fmt.Printf(USAGE, config.Version, config.Author, os.Args[0])
	os.Exit(rc)
}
