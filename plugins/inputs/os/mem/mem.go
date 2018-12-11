package mem

import (
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/lists/loader"
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"

	"github.com/shirou/gopsutil/mem"
)

type InputOSMem struct {}

func (l *InputOSMem) Collect() {
	if ! config.File.OS.Inputs.Mem {
		return
	}

	vmStat, err := mem.VirtualMemory()

	if err == nil {
		metrics.Load().Add(metrics.Metric{
			Key: "zenit_os",
			Tags: []metrics.Tag{
				{"name", "mem"},
			},
			Values: vmStat.UsedPercent,
		})
	}
}

func init() {
	loader.Add("InputOSMem", func() loader.Plugin { return &InputOSMem{} })
}
