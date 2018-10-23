package connections

import (
	"fmt"

	"github.com/swapbyt3s/zenit/common"
	"github.com/swapbyt3s/zenit/common/log"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"
	"github.com/swapbyt3s/zenit/plugins/lists/alerts"
)

func Register() {
	if ! config.File.MySQL.Inputs.Variables {
		log.Info("Require to enable MySQL Variables in config file.")
		return
	}

	if ! config.File.MySQL.Inputs.Status {
		log.Info("Require to enable MySQL Status in config file.")
		return
	}

	var m = metrics.Load()
	var value interface{}
	value = m.FetchOne("zenit_mysql_variables", "name", "max_connections")
	var MaxConnections = float64(common.InterfaceToInt(value))
	value = m.FetchOne("zenit_mysql_status", "name", "Threads_connected")
	var ThreadsConnected = float64(common.InterfaceToInt(value))
	var percentage = int(common.Percentage(ThreadsConnected, MaxConnections))

	if percentage == -1 {
		return
	}

	// Build one message with details for notification:
	var message = fmt.Sprintf("*Current:* %d%%", percentage)

	alerts.Load().Register(
		"connections",
		"MySQL Connections",
		config.File.MySQL.Alerts.Connections.Duration,
		config.File.MySQL.Alerts.Connections.Warning,
		config.File.MySQL.Alerts.Connections.Critical,
		percentage,
		message,
	)
}
