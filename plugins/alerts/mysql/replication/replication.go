package replication

import (
	"fmt"

	"github.com/swapbyt3s/zenit/common"
	"github.com/swapbyt3s/zenit/common/log"
	"github.com/swapbyt3s/zenit/common/mysql"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/lists/accumulator"
	"github.com/swapbyt3s/zenit/plugins/lists/alerts"
)

func Check() {
	if ! config.File.MySQL.Inputs.Slave {
		log.Info("Require to enable MySQL Slave Status in config file.")
		return
	}

	var message string = ""
	var running int
	var value interface{}

	var metrics = accumulator.Load()
	var check = alerts.Load().Exist("replication")

	value = metrics.FetchOne("mysql_slave", "name", "Slave_IO_Running")
	var ioRunning = common.InterfaceToInt(value)
	value = metrics.FetchOne("mysql_slave", "name", "Slave_SQL_Running")
	var sqlRunning = common.InterfaceToInt(value)
	value = metrics.FetchOne("mysql_slave", "name", "Last_SQL_Errno")
	var sqlError = common.InterfaceToInt(value)

	message += fmt.Sprintf("*IO Running:* %s\n", mysql.YesOrNo(ioRunning))
	message += fmt.Sprintf("*SQL Running:* %s\n", mysql.YesOrNo(sqlRunning))
	message += fmt.Sprintf("*SQL Error:* %d\n", sqlError)

	running = 2 - (ioRunning + sqlRunning)

	if check == nil {
		alerts.Load().Add(
			"replication",
			"MySQL Replication Status",
			config.File.MySQL.Alerts.Replication.Duration,
			1, // Warning
			1, // Critical
			running,
			message,
			true,
		)
	} else {
		check.Update(running, message)
	}
}
