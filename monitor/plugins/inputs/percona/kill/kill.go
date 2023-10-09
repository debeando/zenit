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
			log.ErrorWithFields(name, log.Fields{"error": err})
		}
	}()

	if !cnf.Inputs.Process.PerconaToolKitKill {
		return
	}

	var pid = exec.PGrep("pt-kill")
	var value = 0

	if pid > 0 {
		value = 1
	}

	log.DebugWithFields(name, log.Fields{"pt_kill": value})

	mtc.Add(metrics.Metric{
		Key: "process_pt_kill",
		Tags: []metrics.Tag{
			{Name: "hostname", Value: cnf.General.Hostname},
		},
		Values: []metrics.Value{
			{Key: "pt_kill", Value: value},
		},
	})
}

func init() {
	inputs.Add("InputsPerconaKill", func() inputs.Input { return &Plugin{} })
}
