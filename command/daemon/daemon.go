package daemon

import (
	"github.com/swapbyt3s/zenit/common/log"
	"github.com/swapbyt3s/zenit/plugins"

	"github.com/kardianos/service"
)

var (
	srv service.Service
	err error
)

type program struct{}

func Configure() {
	prg := &program{}

	svcConfig := &service.Config{
		Name:        "zenit",
		DisplayName: "Zenit",
		Description: "Zenit Agent",
		Executable:  "/usr/bin/zenit",
	}

	srv, err = service.New(prg, svcConfig)
	if err != nil {
		log.Error("Daemon configure", map[string]interface{}{"error": err})
	}
}

func (p *program) Start(s service.Service) error {
	plugins.Load()

	return nil
}

func (p *program) Stop(s service.Service) error {
	return nil
}

func Daemonize() {
	if err = srv.Run(); err != nil {
		log.Error("Daemon run", map[string]interface{}{"error": err})
	}
}

func Install() {
	if err = srv.Install(); err != nil {
		log.Error("Daemon install", map[string]interface{}{"error": err})
	}
}

func Uninstall() {
	if err = srv.Uninstall(); err != nil {
		log.Error("Daemon uninstall", map[string]interface{}{"error": err})
	}
}
