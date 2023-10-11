package service

import (
	"github.com/debeando/go-common/file"
	"github.com/debeando/go-common/log"
	"github.com/kardianos/service"
	"github.com/spf13/cobra"
)

var (
	install   bool
	uninstall bool
)

type program struct{}

func NewCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "service [OPTIONS]",
		Short: "Run monitoring in daemon mode",
		Example: `
  Install init script:
    $ sudo zenit service --install

  Start zenit agent:
    $ sudo [systemctl|initctl] [start|stop|restart|status] zenit`,
		Run: func(cmd *cobra.Command, args []string) {
			if !install && !uninstall {
				cmd.Help()
				return
			}

			prg := &program{}

			svcConfig := &service.Config{
				Name:        "zenit",
				DisplayName: "Zenit",
				Description: "Zenit Agent",
				Executable:  "/usr/bin/zenit",
				Arguments:   []string{"agent", "--config=/etc/zenit/zenit.yaml"},
			}

			srv, err := service.New(prg, svcConfig)
			if err != nil {
				log.Error(err.Error())
			}

			if install {
				if ! file.Exist("/usr/bin/zenit") {
					log.ErrorWithFields("The executable file was not found.", log.Fields{"path": "/usr/bin/zenit"})
					return
				}

				if err = srv.Install(); err != nil {
					log.Error(err.Error())
				}
			}

			if uninstall {
				if err = srv.Uninstall(); err != nil {
					log.Error(err.Error())
				}
			}
		},
	}

	cmd.Flags().BoolVar(&install, "install", false, "Install service")
	cmd.Flags().BoolVar(&uninstall, "uninstall", false, "Uninstall service")

	return cmd
}

func (p *program) Start(s service.Service) error {
	return nil
}

func (p *program) Stop(s service.Service) error {
	return nil
}
