package disk

import (
	"fmt"
	"strings"

	"zenit/config"
	"zenit/agent/plugins/inputs"
	"zenit/agent/plugins/lists/metrics"

	"github.com/debeando/go-common/log"
	"github.com/shirou/gopsutil/disk"
)

type Plugin struct{}

func (p *Plugin) Collect(name string, cnf *config.Config, mtc *metrics.Items) {
	defer func() {
		if err := recover(); err != nil {
			log.ErrorWithFields(name, log.Fields{"message": err})
		}
	}()

	if !cnf.Inputs.OS.Disk {
		return
	}

	devices, err := disk.Partitions(false)
	if err != nil {
		return
	}

	var v = metrics.Values{}

	for _, device := range devices {
		u, err := disk.Usage(device.Mountpoint)

		if err != nil {
			return
		}

		log.DebugWithFields(name, log.Fields{"device": GetDevice(device.Device), "usage": fmt.Sprintf("%.2f", u.UsedPercent)})

		v.Add(metrics.Value{
			Key:   GetDevice(device.Device),
			Value: u.UsedPercent,
		})
	}

	mtc.Add(metrics.Metric{
		Key: "os_disk",
		Tags: []metrics.Tag{
			{Name: "hostname", Value: cnf.General.Hostname},
		},
		Values: v,
	})
}

func GetDevice(s string) string {
	return strings.Replace(s, "/dev/", "", -1)
}

func init() {
	inputs.Add("InputOSDisk", func() inputs.Input { return &Plugin{} })
}
