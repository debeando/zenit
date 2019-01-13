package disk

import (
	"fmt"
	"strings"

	"github.com/swapbyt3s/zenit/common/log"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"
	"github.com/swapbyt3s/zenit/plugins/inputs"

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

		log.Debug(fmt.Sprintf("Plugin - InputOSDisk - Disk=%s(%.2f)", GetDevice(device.Device), u.UsedPercent))
	}
}

func GetDevice(s string) string {
	return strings.Replace(s, "/dev/", "", -1)
}

func init() {
	inputs.Add("InputOSDisk", func() inputs.Input { return &InputOSDisk{} })
}
