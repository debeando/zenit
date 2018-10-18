package os

import (
	"github.com/swapbyt3s/zenit/plugins/lists/accumulator"

	"github.com/shirou/gopsutil/disk"
)

func Disk() {
	devices, err := disk.Partitions(false)

	if err != nil {
		return
	}

	for _, device := range devices {
		u, err := disk.Usage(device.Mountpoint)

		if err != nil {
			return
		}

		accumulator.Load().Add(accumulator.Metric{
			Key: "os",
			Tags: []accumulator.Tag{
				{"name", "disk"},
				{"device", device.Device},
			},
			Values: u.UsedPercent,
		})
	}
}
