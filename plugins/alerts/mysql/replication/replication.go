package replication

import (
	"fmt"

	"github.com/swapbyt3s/zenit/common"
	"github.com/swapbyt3s/zenit/common/log"
	"github.com/swapbyt3s/zenit/common/mysql"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"
	"github.com/swapbyt3s/zenit/plugins/lists/alerts"
)

func Register() {
	if ! config.File.MySQL.Inputs.Slave {
		log.Info("Require to enable MySQL Slave Status in config file.")
		return
	}

	var message string = ""
	var running int
	var value interface{}

	var metrics = metrics.Load()

	value = metrics.FetchOne("mysql_slave", "name", "Slave_IO_Running")
	var ioRunning = common.InterfaceToInt(value)
	value = metrics.FetchOne("mysql_slave", "name", "Slave_SQL_Running")
	var sqlRunning = common.InterfaceToInt(value)
	value = metrics.FetchOne("mysql_slave", "name", "Last_SQL_Errno")
	var sqlError = common.InterfaceToInt(value)

	if sqlError == -1 {
		return
	}

	message += fmt.Sprintf("*IO Running:* %s\n", mysql.YesOrNo(ioRunning))
	message += fmt.Sprintf("*SQL Running:* %s\n", mysql.YesOrNo(sqlRunning))
	message += fmt.Sprintf("*SQL Error:* %d\n", sqlError)

	running = 2 - (ioRunning + sqlRunning)

	alerts.Load().Register(
		"replication",
		"MySQL Replication Status",
		config.File.MySQL.Alerts.Replication.Duration,
		1, // Warning
		1, // Critical
		running,
		message,
	)
}
