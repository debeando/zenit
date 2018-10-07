package os

import (
	"github.com/swapbyt3s/zenit/plugins/accumulator"

	"github.com/shirou/gopsutil/cpu"
)

func CPU() {
	percentage, err := cpu.Percent(0, false)

	if err == nil {
		accumulator.Load().Add(accumulator.Metric{
			Key: "os",
			Tags: []accumulator.Tag{
				{"name", "cpu"},
			},
			Values: percentage[0],
		})
	}
}
