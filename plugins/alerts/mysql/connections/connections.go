package connections

import (
	"fmt"

	"github.com/swapbyt3s/zenit/common"
	"github.com/swapbyt3s/zenit/common/log"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/alerts"
	"github.com/swapbyt3s/zenit/plugins/lists/checks"
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"
)

type MySQLConnections struct {}

func (l *MySQLConnections) Collect() {
	if ! config.File.MySQL.Alerts.Connections.Enable {
		return
	}

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
	var MaxConnections = common.InterfaceToUInt64(value)
	value = m.FetchOne("zenit_mysql_status", "name", "Threads_connected")
	var ThreadsConnected = common.InterfaceToUInt64(value)
	var percentage = uint64(common.Percentage(ThreadsConnected, MaxConnections))

	//if MaxConnections == 0 || ThreadsConnected == 0 || percentage == 0 {
	//	return
	//}

	// Build one message with details for notification:
	var message = fmt.Sprintf("*Current:* %d%%", percentage)

	checks.Load().Register(
		"connections",
		"MySQL Connections",
		config.File.MySQL.Alerts.Connections.Duration,
		config.File.MySQL.Alerts.Connections.Warning,
		config.File.MySQL.Alerts.Connections.Critical,
		percentage,
		message,
	)

	log.Debug(
		fmt.Sprintf("Plugin - AlertMySQLConnections - MaxConnections: %d, ThreadsConnected: %d, Percentage: %d%%",
			MaxConnections,
			ThreadsConnected,
			percentage,
		),
	)
}

func init() {
	alerts.Add("AlertMySQLConnections", func() alerts.Alert { return &MySQLConnections{} })
}
