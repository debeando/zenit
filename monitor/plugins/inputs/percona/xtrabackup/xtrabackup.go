package process

import (
	"zenit/config"
	"zenit/monitor/plugins/inputs"
	"zenit/monitor/plugins/lists/metrics"

	"github.com/debeando/go-common/exec"
	"github.com/debeando/go-common/log"
)

type Plugin struct{}

func (p *Plugin) Collect(name string, cnf *config.Config, mtc *metrics.Items) {
	defer func() {
		if err := recover(); err != nil {
			log.ErrorWithFields(name, log.Fields{"message": err})
		}
	}()

	if !cnf.Inputs.Process.PerconaXtraBackup {
		return
	}

	var pid = exec.PGrep("xtrabackup")
	var value = 0

	if pid > 0 {
		value = 1
	}

	log.DebugWithFields(name, log.Fields{"xtrabackup": value})

	mtc.Add(metrics.Metric{
		Key: "process_xtrabackup",
		Tags: []metrics.Tag{
			{Name: "hostname", Value: cnf.General.Hostname},
		},
		Values: []metrics.Value{
			{Key: "xtrabackup", Value: value},
		},
	})
}

func init() {
	inputs.Add("InputsPerconaXtraBackup", func() inputs.Input { return &Plugin{} })
}
