package disk

import (
	"strings"

	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/lists/loader"
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"

	"github.com/shirou/gopsutil/disk"
)

type InputOSDisk struct {}

func (l *InputOSDisk) Collect() {
	if ! config.File.OS.Inputs.Disk {
		return
	}

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

func init() {
	loader.Add("InputOSDisk", func() loader.Plugin { return &InputOSDisk{} })
}
