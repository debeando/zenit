package cpu

import (
	"fmt"
	"log"

	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/accumulator"
	"github.com/swapbyt3s/zenit/plugins/lists/alerts"
)

func Check() {
	if ! config.File.OS.Alerts.CPU.Enable {
		log.Printf("W! - Require to enable OS CPU in config file.\n")
		return
	}

	var metrics = accumulator.Load()
	var check = alerts.Load().Exist("cpu")
	var message string = ""
	var percentage1 = metrics.Find("os", "name", "cpu")

	if percentage2, ok2 := percentage1.(float64); ok2 {
		if ! ok2 {
			return
		}

		percentage3 := int(percentage2)

		message += fmt.Sprintf("*CPU:* %d\n", percentage3)

		if check == nil {
			log.Printf("D! - Alert:OS:CPU - Adding\n")
			log.Printf("D! - Alert:OS:CPU - Message: %s\n", message)
			alerts.Load().Add(
				"cpu",
				"CPU",
				config.File.OS.Alerts.CPU.Duration,
				config.File.OS.Alerts.CPU.Warning,
				config.File.OS.Alerts.CPU.Critical,
				percentage3,
				message,
				true,
			)
		} else {
			log.Printf("D! - Alert:OS:CPU - Message: %s\n", message)
			log.Printf("D! - Alert:OS:CPU - Updateing\n")
			check.Update(percentage3, message)
		}
	}
}

func Float64ToInt(value interface{}) int {
	if v, ok := value.(float64); ok {
		return int(v)
	}
	return -1
}
