package disk

import (
	"fmt"
	"strings"

	"github.com/swapbyt3s/zenit/common/log"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/inputs"
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"

	"github.com/shirou/gopsutil/disk"
)

type InputOSDisk struct{}

func (l *InputOSDisk) Collect() {
	defer func() {
		if err := recover(); err != nil {
			log.Debug(fmt.Sprintf("Plugin - InputOSDisk - Panic (code %d) has been recover from somewhere.\n", err))
		}
	}()

	if !config.File.Inputs.OS.Disk {
		return
	}

	devices, err := disk.Partitions(false)
	if err != nil {
		return
	}

	var a = metrics.Load()
	var v = []metrics.Value{}

	for _, device := range devices {
		u, err := disk.Usage(device.Mountpoint)

		if err != nil {
			return
		}

		log.Debug(fmt.Sprintf("Plugin - InputOSDisk - Disk=%s(%.2f)", GetDevice(device.Device), u.UsedPercent))

		v = append(v, metrics.Value{
			Key: GetDevice(device.Device),
			Value: u.UsedPercent,
		})
	}

	a.Add(metrics.Metric{
		Key:    "os",
		Tags:   []metrics.Tag{},
		Values: v,
	})
}

func GetDevice(s string) string {
	return strings.Replace(s, "/dev/", "", -1)
}

func init() {
	inputs.Add("InputOSDisk", func() inputs.Input { return &InputOSDisk{} })
}
