package cpu

import (
	"fmt"

	"github.com/debeando/zenit/common/log"
	"github.com/debeando/zenit/config"
	"github.com/debeando/zenit/plugins/inputs"
	"github.com/debeando/zenit/plugins/lists/metrics"

	"github.com/shirou/gopsutil/cpu"
)

type InputOSCPU struct{}

func (l *InputOSCPU) Collect() {
	defer func() {
		if err := recover(); err != nil {
			log.Error("InputOSCPU", map[string]interface{}{"error": err})
		}
	}()

	if !config.File.Inputs.OS.CPU {
		return
	}

	var a = metrics.Load()

	percentage, err := cpu.Percent(0, false)
	if err == nil {
		a.Add(metrics.Metric{
			Key: "os_cpu",
			Tags: []metrics.Tag{
				{Name: "hostname", Value: config.File.General.Hostname},
			},
			Values: []metrics.Value{
				{Key: "percentage", Value: percentage[0]},
			},
		})
	}

	log.Debug("InputOSCPU", map[string]interface{}{"value": fmt.Sprintf("%.2f", percentage[0])})
}

func init() {
	inputs.Add("InputOSCPU", func() inputs.Input { return &InputOSCPU{} })
}
