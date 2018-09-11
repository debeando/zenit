package os

import (
	"strings"
	"time"

	"github.com/swapbyt3s/zenit/common"
	"github.com/swapbyt3s/zenit/common/file"
	"github.com/swapbyt3s/zenit/plugins/accumulator"
)

type CPUMetric struct {
	idle  uint64
	total uint64
}

func CPU() {
	c := [2]CPUMetric{}

	c[0].idle, c[0].total = getCPUSample()
	time.Sleep(300 * time.Millisecond)
	c[1].idle, c[1].total = getCPUSample()

	if c[0].total > 0 && c[1].total > 0 && c[0].idle > 0 && c[1].idle > 0 {
		total := c[1].total - c[0].total
		idle := c[1].idle - c[0].idle
		percentage := (float64(total) - float64(idle)) / float64(total) * 100.0

		accumulator.Load().Add(accumulator.Metric{
			Key: "os",
			Tags: []accumulator.Tag{{"system", "linux"},
				{"hardware", "cpu"}},
			Values: []accumulator.Value{{"used_percent", percentage}},
		})
	}
}

func getCPUSample() (idle uint64, total uint64) {
	lines := file.Read("/proc/stat")
	if len(lines) > 0 {
		fields := strings.Fields(lines)

		for i := 1; i < len(fields); i++ {
			total += common.StringToUInt64(fields[i])
			if i == 4 {
				idle = common.StringToUInt64(fields[i])
			}
		}
	}
	return
}
