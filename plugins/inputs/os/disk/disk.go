package disk

import (
	"strings"

	"github.com/swapbyt3s/zenit/plugins/lists/metrics"

	"github.com/shirou/gopsutil/disk"
)

func Collect() {
	devices, err := disk.Partitions(false)

	if err != nil {
		return
	}

	for _, device := range devices {
		u, err := disk.Usage(device.Mountpoint)

		if err != nil {
			return
		}

		metrics.Load().Add(metrics.Metric{
			Key: "zenit_os",
			Tags: []metrics.Tag{
				{"name", "disk"},
				{"device", GetDevice(device.Device)},
			},
			Values: u.UsedPercent,
		})
	}
}

func GetDevice(s string) string {
	return strings.Replace(s, "/dev/", "", -1)
}
