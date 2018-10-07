package os

import (
	"github.com/swapbyt3s/zenit/plugins/accumulator"

	"github.com/shirou/gopsutil/disk"
)

func Disk() {
	parts, err := disk.Partitions(false)

	if err != nil {
		return
	}

	for _, part := range parts {
		u, err := disk.Usage(part.Mountpoint)

		if err != nil {
			return
		}

		accumulator.Load().Add(accumulator.Metric{
			Key: "os",
			Tags: []accumulator.Tag{
				{"name", "disk"},
				{"filesystem", part.Fstype},
				{"device", part.Device},
				{"mount", part.Mountpoint},
			},
			Values: u.UsedPercent,
		})
	}
}
