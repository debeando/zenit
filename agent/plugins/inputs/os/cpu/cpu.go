package cpu

import (
	"fmt"

	"zenit/agent/plugins/inputs"
	"zenit/agent/plugins/lists/metrics"
	"zenit/config"

	"github.com/debeando/go-common/log"
	"github.com/shirou/gopsutil/cpu"
)

type Plugin struct{}

func (p *Plugin) Collect(name string, cnf *config.Config, mtc *metrics.Items) {
	defer func() {
		if err := recover(); err != nil {
			log.ErrorWithFields(name, log.Fields{"message": err})
		}
	}()

	if !cnf.Inputs.OS.CPU {
		return
	}

	percentage, err := cpu.Percent(0, false)
	if err == nil {
		mtc.Add(metrics.Metric{
			Key: "os_cpu",
			Tags: []metrics.Tag{
				{Name: "hostname", Value: cnf.General.Hostname},
			},
			Values: []metrics.Value{
				{Key: "percentage", Value: percentage[0]},
			},
		})
	}

	log.DebugWithFields(name, log.Fields{"value": fmt.Sprintf("%.2f", percentage[0])})
}

func init() {
	inputs.Add("InputOSCPU", func() inputs.Input { return &Plugin{} })
}
