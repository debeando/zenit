package readonly

import (
	"fmt"
	"log"

	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/accumulator"
	"github.com/swapbyt3s/zenit/plugins/lists/alerts"
)

func Check() {
	if ! config.File.MySQL.Variables {
		log.Printf("W! - Require to enable MySQL Variables in config file.\n")
		return
	}

	var metrics = accumulator.Load()
	var value, ok = metrics.Find("mysql_variables", "read_only")

	// Verify find match:
	if ! ok {
		return
	}

	// Verify value range is valid:
	if ! (value == 0 || value == 1) {
		return
	}

	// Invert value for compatibility alert evaluation:
	value = (value ^ 1)

	// Find own check:
	var check = alerts.Load().Exist("readonly")

	// Build one message with details for notification:
	var message = fmt.Sprintf("*Current:* %b", (value ^ 1))

	// Register new check and update last status:
	if check == nil {
		alerts.Load().Add(
			"readonly",
			"MySQL Read Only",
			config.File.MySQL.Alerts.ReadOnly.Duration,
			1, // Warning
			1, // Critical
			value,
			message,
			true,
		)
	} else {
		check.Update(value, message)
	}
}
