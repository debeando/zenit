package lagging

import (
	"fmt"

	"github.com/swapbyt3s/zenit/common"
	"github.com/swapbyt3s/zenit/common/log"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/alerts"
	"github.com/swapbyt3s/zenit/plugins/lists/checks"
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"
)

type MySQLLagging struct {}

func (l *MySQLLagging) Collect() {
	if ! config.File.MySQL.Alerts.Replication.Enable {
		return
	}

	if ! config.File.MySQL.Inputs.Slave {
		log.Info("Require to enable MySQL Slave Status in config file.")
		return
	}

	var metrics = metrics.Load()
	var value = metrics.FetchOne("zenit_mysql_slave", "name", "Seconds_Behind_Master")
	var lagging = common.InterfaceToUInt64(value)

	// Build one message with details for notification:
	var message = fmt.Sprintf("*Lagging:* %d\n", lagging)

	checks.Load().Register(
		"lagging",
		"MySQL Replication Lagging",
		config.File.MySQL.Alerts.Replication.Duration,
		config.File.MySQL.Alerts.Replication.Warning,
		config.File.MySQL.Alerts.Replication.Critical,
		lagging,
		message,
	)
}

func init() {
	alerts.Add("AlertMySQLLagging", func() alerts.Alert { return &MySQLLagging{} })
}
