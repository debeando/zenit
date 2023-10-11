package process

import (
	"zenit/agent/plugins/inputs"
	"zenit/agent/plugins/lists/metrics"
	"zenit/config"

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

	if !cnf.Inputs.Process.PerconaToolKitSlaveDelay {
		return
	}

	var pid = exec.PGrep("pt-slave-delay")
	var value int64 = 0

	if pid > 0 {
		value = 1
	}

	log.DebugWithFields(name, log.Fields{"pt_slave_delay": value})

	mtc.Add(metrics.Metric{
		Key: "process_pt_slave_delay",
		Tags: []metrics.Tag{
			{Name: "hostname", Value: cnf.General.Hostname},
		},
		Values: []metrics.Value{
			{Key: "pt_slave_delay", Value: value},
		},
	})
}

func init() {
	inputs.Add("InputsPerconaToolkitSlaveDelay", func() inputs.Input { return &Plugin{} })
}
