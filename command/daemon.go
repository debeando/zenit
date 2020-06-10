package command

import (
	"github.com/swapbyt3s/zenit/common/log"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins"

	"github.com/kardianos/service"
)

var (
	daemon service.Service
	err error
)

type program struct{}

func (p *program) Start(s service.Service) error {
	config.File.Load()
	config.File.SanityCheck()
	plugins.Load()

	return nil
}

func (p *program) Stop(s service.Service) error {
	return nil
}

func init() {
	prg := &program{}

	svcConfig := &service.Config{
		Name:        "zenit",
		DisplayName: "Zenit",
		Description: "Zenit Agent",
		Executable:  "/usr/bin/zenit",
	}

	daemon, err = service.New(prg, svcConfig)
	if err != nil {
		log.Error("Daemon", map[string]interface{}{"error": err})
	}
}

func Daemonize() {
	if err = daemon.Run(); err != nil {
		log.Error("Daemon", map[string]interface{}{"error": err})
	}
}

func Install() {
	if err = daemon.Install(); err != nil {
		log.Error("Daemon", map[string]interface{}{"error": err})
	}
}

func Uninstall() {
	if err = daemon.Uninstall(); err != nil {
		log.Error("Daemon", map[string]interface{}{"error": err})
	}
}
