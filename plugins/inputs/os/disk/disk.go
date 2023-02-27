package disk

import (
	"fmt"
	"strings"

	"github.com/debeando/zenit/common/log"
	"github.com/debeando/zenit/config"
	"github.com/debeando/zenit/plugins/inputs"
	"github.com/debeando/zenit/plugins/lists/metrics"

	"github.com/shirou/gopsutil/disk"
)

type InputOSDisk struct{}

func (l *InputOSDisk) Collect() {
	defer func() {
		if err := recover(); err != nil {
			log.Error("InputOSDisk", map[string]interface{}{"error": err})
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

		log.Debug("InputOSDisk", map[string]interface{}{"device": GetDevice(device.Device), "usage": fmt.Sprintf("%.2f", u.UsedPercent)})

		v = append(v, metrics.Value{
			Key:   GetDevice(device.Device),
			Value: u.UsedPercent,
		})
	}

	a.Add(metrics.Metric{
		Key: "os_disk",
		Tags: []metrics.Tag{
			{"hostname", config.File.General.Hostname},
		},
		Values: v,
	})
}

func GetDevice(s string) string {
	return strings.Replace(s, "/dev/", "", -1)
}

func init() {
	inputs.Add("InputOSDisk", func() inputs.Input { return &InputOSDisk{} })
}
