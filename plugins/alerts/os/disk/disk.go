package disk

import (
	"fmt"
	"log"

	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/lists/accumulator"
	"github.com/swapbyt3s/zenit/plugins/lists/alerts"
)

func Check() {
	if ! config.File.OS.Alerts.Disk.Enable {
		log.Printf("W! - Require to enable OS Disk in config file.\n")
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

					var check = alerts.Load().Exist("disk_" + device)
					var percentage = Float64ToInt(metric.Values)

					message += fmt.Sprintf("*Volume:* %s, *Usage:* %d%%\n", device, percentage)

					if check == nil {
						log.Printf("D! - Alert:OS:Disk - Adding\n")
						log.Printf("D! - Alert:OS:Disk - Message: %s\n", message)
						alerts.Load().Add(
							"disk_" + device,
							"Volumen",
							config.File.OS.Alerts.Disk.Duration,
							config.File.OS.Alerts.Disk.Warning,
							config.File.OS.Alerts.Disk.Critical,
							percentage,
							message,
							true,
						)
					} else {
						log.Printf("D! - Alert:OS:Disk - Message: %s\n", message)
						log.Printf("D! - Alert:OS:Disk - Updateing\n")
						check.Update(percentage, message)
					}
				}
			}
		}
	}
}

func Float64ToInt(value interface{}) int {
	if v, ok := value.(float64); ok {
		return int(v)
	}
	return -1
}
