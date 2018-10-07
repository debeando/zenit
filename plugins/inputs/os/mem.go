package os

import (
	"github.com/swapbyt3s/zenit/plugins/accumulator"

	"github.com/shirou/gopsutil/mem"
)

func Mem() {
	vmStat, err := mem.VirtualMemory()

	if err == nil {
		accumulator.Load().Add(accumulator.Metric{
			Key: "os",
			Tags: []accumulator.Tag{
				{"name", "mem"},
			},
			Values: vmStat.UsedPercent,
		})
	}
}
