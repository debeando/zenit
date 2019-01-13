package replication

import (
	"fmt"

	"github.com/swapbyt3s/zenit/common"
	"github.com/swapbyt3s/zenit/common/log"
	"github.com/swapbyt3s/zenit/common/mysql"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/alerts"
	"github.com/swapbyt3s/zenit/plugins/lists/checks"
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"
)

type MySQLReplication struct {}

func (l *MySQLReplication) Collect() {
	if ! config.File.MySQL.Alerts.Replication.Enable {
		return
	}

	if ! config.File.MySQL.Inputs.Slave {
		log.Info("Require to enable MySQL Slave Status in config file.")
		return
	}

	var message string = ""
	var running uint64
	var value interface{}

	var metrics = metrics.Load()

	value = metrics.FetchOne("zenit_mysql_slave", "name", "Slave_IO_Running")
	var ioRunning = common.InterfaceToUInt64(value)
	value = metrics.FetchOne("zenit_mysql_slave", "name", "Slave_SQL_Running")
	var sqlRunning = common.InterfaceToUInt64(value)
	value = metrics.FetchOne("zenit_mysql_slave", "name", "Last_SQL_Errno")
	var sqlError = common.InterfaceToUInt64(value)

	if sqlRunning == 0 {
		var delay = metrics.FetchOne("zenit_process", "name", "pt_slave_delay")

//		if delay == -1 {
//			return
//		}

		if delay == 1 {
			log.Debug("Plugin - AlertMySQLReplication - Is running pt-slave-delay, skip check!.")

			return
		}
	}

//	if sqlError == -1 {
//		return
//	}

	message += fmt.Sprintf("*IO Running:* %s\n", mysql.YesOrNo(ioRunning))
	message += fmt.Sprintf("*SQL Running:* %s\n", mysql.YesOrNo(sqlRunning))
	message += fmt.Sprintf("*SQL Error:* %d\n", sqlError)

	running = 2 - (ioRunning + sqlRunning)

	checks.Load().Register(
		"replication",
		"MySQL Replication Status",
		config.File.MySQL.Alerts.Replication.Duration,
		1, // Warning
		1, // Critical
		running,
		message,
	)

	log.Debug(fmt.Sprintf("Plugin - AlertMySQLReplication - Slave_IO_Running=%d", ioRunning))
	log.Debug(fmt.Sprintf("Plugin - AlertMySQLReplication - Slave_SQL_Running=%d", sqlRunning))
	log.Debug(fmt.Sprintf("Plugin - AlertMySQLReplication - Last_SQL_Errno=%d", sqlError))
}

func init() {
	alerts.Add("AlertMySQLReplication", func() alerts.Alert { return &MySQLReplication{} })
}
