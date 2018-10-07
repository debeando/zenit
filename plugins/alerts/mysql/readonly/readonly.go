package readonly

import (
	"fmt"
	"log"

	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/accumulator"
	"github.com/swapbyt3s/zenit/plugins/lists/alerts"
)

func Check() {
	if ! config.File.MySQL.Inputs.Variables {
		log.Printf("W! - Require to enable MySQL Variables in config file.\n")
		return
	}

	var metrics = accumulator.Load()
	var value = metrics.Find("mysql_variables", "name", "read_only")
	var status = Float64ToInt(value)

	// Verify status range is valid:
	if ! (status == 0 || status == 1) {
		return
	}

	// Invert status for compatibility alert evaluation:
	status = (status ^ 1)

	// Find own check:
	var check = alerts.Load().Exist("readonly")

	// Build one message with details for notification:
	var message = fmt.Sprintf("*Current:* %s", YesOrNo(status ^ 1))

	// Register new check and update last status:
	if check == nil {
		alerts.Load().Add(
			"readonly",
			"MySQL Read Only",
			config.File.MySQL.Alerts.ReadOnly.Duration,
			1, // Warning
			1, // Critical
			status,
			message,
			true,
		)
	} else {
		check.Update(status, message)
	}
}

func YesOrNo(v int) string {
	if v == 1 {
		return "Yes"
	}
	return "No"
}

func Float64ToInt(value interface{}) int {
	if v, ok := value.(float64); ok {
		return int(v)
	}
	return -1
}
