package cpu

import (
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"

	"github.com/shirou/gopsutil/cpu"
)

func Collect() {
	percentage, err := cpu.Percent(0, false)

	if err == nil {
		metrics.Load().Add(metrics.Metric{
			Key: "os",
			Tags: []metrics.Tag{
				{"name", "cpu"},
			},
			Values: percentage[0],
		})
	}
}
