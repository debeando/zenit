package os

import (
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"

	"github.com/shirou/gopsutil/mem"
)

func Mem() {
	vmStat, err := mem.VirtualMemory()

	if err == nil {
		metrics.Load().Add(metrics.Metric{
			Key: "os",
			Tags: []metrics.Tag{
				{"name", "mem"},
			},
			Values: vmStat.UsedPercent,
		})
	}
}
