package disk

import (
	"fmt"

	"github.com/swapbyt3s/zenit/common"
	"github.com/swapbyt3s/zenit/common/log"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/lists/accumulator"
	"github.com/swapbyt3s/zenit/plugins/lists/alerts"
)

func Check() {
	if ! config.File.OS.Alerts.Disk.Enable {
		log.Info("Require to enable OS Disk in config file.")
		return
	}

	var metrics = accumulator.Load()

	for _, metric := range *metrics {
		if metric.Key == "os" {
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

					// fmt.Printf("--> DEBUG! value before: %#v\n", metric.Values)
					var percentage = int(common.InterfaceToFloat64(metric.Values))
					// fmt.Printf("--> DEBUG! value after: %#v\n", percentage)

					message += fmt.Sprintf("*Volume:* %s, *Usage:* %d%%\n", device, percentage)

					alerts.Load().Register(
						"disk_" + device,
						"Volumen",
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
