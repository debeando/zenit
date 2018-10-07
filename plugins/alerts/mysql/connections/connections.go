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
	var value interface{}
	value = metrics.Find("mysql_variables", "name", "max_connections")
	var MaxConnections = InterfaceToFloat64(value)
	value = metrics.Find("mysql_status", "name", "Threads_connected")
	var ThreadsConnected = InterfaceToFloat64(value)
	var percentage = int(common.Percentage(ThreadsConnected, MaxConnections))
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
			percentage,
			message,
			true,
		)
	} else {
		check.Update(percentage, message)
	}
}

func InterfaceToFloat64(value interface{}) float64 {
	if v, ok := value.(float64); ok {
		return v
	}
	return -1
}
