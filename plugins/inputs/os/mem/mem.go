package mem

import (
	"fmt"

	"github.com/swapbyt3s/zenit/common/log"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/inputs"
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"

	"github.com/shirou/gopsutil/mem"
)

type InputOSMem struct{}

func (l *InputOSMem) Collect() {
	defer func() {
		if err := recover(); err != nil {
			log.Debug(fmt.Sprintf("Plugin - InputOSMem - Panic (code %d) has been recover from somewhere.\n", err))
		}
	}()

	if !config.File.Inputs.OS.Mem {
		return
	}

	var a = metrics.Load()

	vmStat, err := mem.VirtualMemory()
	if err == nil {
		a.Add(metrics.Metric{
			Key: "os",
			Tags: []metrics.Tag{
				{"hostname", config.File.General.Hostname},
			},
			Values: []metrics.Value{
				{ "mem", vmStat.UsedPercent },
			},
		})

		log.Debug(fmt.Sprintf("Plugin - InputOSMem - MEM=%.2f", vmStat.UsedPercent))
	}
}

func init() {
	inputs.Add("InputOSMem", func() inputs.Input { return &InputOSMem{} })
}
