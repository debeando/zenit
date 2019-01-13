package disk

import (
	"fmt"

	"github.com/swapbyt3s/zenit/common"
	"github.com/swapbyt3s/zenit/common/log"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/alerts"
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"
	"github.com/swapbyt3s/zenit/plugins/lists/checks"
)

type OSDisk struct {}

func (l *OSDisk) Collect() {
	if ! config.File.OS.Alerts.Disk.Enable {
		log.Info("Require to enable OS Disk in config file.")
		return
	}

	var metrics = metrics.Load()

	for _, metric := range *metrics {
		if metric.Key == "zenit_os" {
			for _, metricTag := range metric.Tags {
				if metricTag.Name == "name" && metricTag.Value == "disk" {
					var message string = ""
					var device string

					for _, tag := range metric.Tags {
						if tag.Name == "device" {
							device = tag.Value
							break
						}
					}

					var percentage = uint64(common.InterfaceToFloat64(metric.Values))

					message += fmt.Sprintf("*Volume:* %s, *Usage:* %d%%\n", device, percentage)

					checks.Load().Register(
						"disk_" + device,
						fmt.Sprintf("Volumen (%s)", device),
						config.File.OS.Alerts.Disk.Duration,
						config.File.OS.Alerts.Disk.Warning,
						config.File.OS.Alerts.Disk.Critical,
						percentage,
						message,
					)
				}
			}
		}
	}
}

func init() {
	alerts.Add("AlertOSDisk", func() alerts.Alert { return &OSDisk{} })
}
