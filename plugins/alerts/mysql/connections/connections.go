package connections

import (
	"fmt"
	"log"

	"github.com/swapbyt3s/zenit/common"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/accumulator"
	"github.com/swapbyt3s/zenit/plugins/lists/alerts"
)

func Check() {
	if ! config.File.MySQL.Inputs.Variables {
		log.Printf("W! - Require to enable MySQL Variables in config file.\n")
		return
	}

	if ! config.File.MySQL.Inputs.Status {
		log.Printf("W! - Require to enable MySQL Status in config file.\n")
		return
	}

	var metrics = accumulator.Load()
	var max_connections uint64
	var threads_connected uint64
	var ok bool

	max_connections, ok = metrics.Find("mysql_variables", "max_connections")
	if ! ok {
		return
	}
	threads_connected, ok = metrics.Find("mysql_status", "Threads_connected")
	if ! ok {
		return
	}

	var percentage = common.Percentage(threads_connected, max_connections)
	var value = common.FloatToUInt(percentage)
	var check = alerts.Load().Exist("connections")

	// Build one message with details for notification:
	var message = fmt.Sprintf("*Current:* %d%%", value)

	if check == nil {
		alerts.Load().Add(
			"connections",
			"MySQL Connections",
			config.File.MySQL.Alerts.Connections.Duration,
			config.File.MySQL.Alerts.Connections.Warning,
			config.File.MySQL.Alerts.Connections.Critical,
			value,
			message,
			true,
		)
	} else {
		check.Update(value, message)
	}
}
