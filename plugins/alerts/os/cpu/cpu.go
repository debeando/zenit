package cpu

import (
	"fmt"
	"log"

	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/common"
	"github.com/swapbyt3s/zenit/plugins/lists/accumulator"
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
	var value = metrics.FetchOne("os", "name", "cpu")
	var percentage = common.InterfaceToInt(value)

	message += fmt.Sprintf("*CPU:* %d\n", percentage)

	if check == nil {
		log.Printf("D! - Alert:OS:CPU - Adding\n")
		log.Printf("D! - Alert:OS:CPU - Message: %s\n", message)
		alerts.Load().Add(
			"cpu",
			"CPU",
			config.File.OS.Alerts.CPU.Duration,
			config.File.OS.Alerts.CPU.Warning,
			config.File.OS.Alerts.CPU.Critical,
			percentage,
			message,
			true,
		)
	} else {
		log.Printf("D! - Alert:OS:CPU - Message: %s\n", message)
		log.Printf("D! - Alert:OS:CPU - Updateing\n")
		check.Update(percentage, message)
	}
}
