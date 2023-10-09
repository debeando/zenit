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

	if !cnf.Inputs.Process.PerconaToolKitDeadlockLogger {
		return
	}

	var pid = exec.PGrep("pt-deadlock-logger")
	var value = 0

	if pid > 0 {
		value = 1
	}

	log.DebugWithFields(name, log.Fields{"pt_deadlock_logger": value})

	mtc.Add(metrics.Metric{
		Key: "process_pt_deadlock_logger",
		Tags: []metrics.Tag{
			{Name: "hostname", Value: cnf.General.Hostname},
		},
		Values: []metrics.Value{
			{Key: "pt_deadlock_logger", Value: value},
		},
	})
}

func init() {
	inputs.Add("InputsPerconaDeadlock", func() inputs.Input { return &Plugin{} })
}
