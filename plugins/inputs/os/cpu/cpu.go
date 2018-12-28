package cpu

import (
	"fmt"

	"github.com/swapbyt3s/zenit/common/log"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/lists/loader"
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"

	"github.com/shirou/gopsutil/cpu"
)

type InputOSCPU struct {}

func (l *InputOSCPU) Collect() {
	if ! config.File.OS.Inputs.CPU {
		return
	}

	percentage, err := cpu.Percent(0, false)

	if err == nil {
		metrics.Load().Add(metrics.Metric{
			Key: "zenit_os",
			Tags: []metrics.Tag{
				{"name", "cpu"},
			},
			Values: percentage[0],
		})
	}

	log.Debug(fmt.Sprintf("Plugin - InputOSCPU - CPU=%.2f", percentage))
}

func init() {
	loader.Add("InputOSCPU", func() loader.Plugin { return &InputOSCPU{} })
}
