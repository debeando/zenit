package lagging

import (
	"fmt"

	"github.com/swapbyt3s/zenit/common"
	"github.com/swapbyt3s/zenit/common/log"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/lists/alerts"
	"github.com/swapbyt3s/zenit/plugins/lists/loader"
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"
)

type MySQLLagging struct {}

func (l *MySQLLagging) Collect() {
	if ! config.File.MySQL.Inputs.Slave {
		log.Info("Require to enable MySQL Slave Status in config file.")
		return
	}

	var metrics = metrics.Load()
	var value = metrics.FetchOne("zenit_mysql_slave", "name", "Seconds_Behind_Master")
	var lagging = common.InterfaceToInt(value)

	if lagging == -1 {
		return
	}

	// Build one message with details for notification:
	var message = fmt.Sprintf("*Lagging:* %d\n", lagging)

	alerts.Load().Register(
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
	loader.Add("AlertMySQLLagging", func() loader.Plugin { return &MySQLLagging{} })
}
