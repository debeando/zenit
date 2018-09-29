package replication

import (
	"fmt"
	"log"

	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/accumulator"
	"github.com/swapbyt3s/zenit/plugins/lists/alerts"
)

func Check() {
	if ! config.File.MySQL.Inputs.Slave {
		log.Printf("W! - Require to enable MySQL Slave Status in config file.\n")
		return
	}

	var ioRunning uint64
	var message string = ""
	var ok bool
	var delay uint64
	var sqlError uint64
	var sqlRunning uint64
	var decide bool

	var metrics = accumulator.Load()
	var check = alerts.Load().Exist("replication")

	ioRunning, ok = metrics.Find("mysql_slave", "Slave_IO_Running")
	if ! ok {
		return
	}
	sqlRunning, ok = metrics.Find("mysql_slave", "Slave_SQL_Running")
	if ! ok {
		return
	}
	sqlError, ok = metrics.Find("mysql_slave", "Last_SQL_Errno")
	if ! ok {
		return
	}
	delay, _ = metrics.Find("mysql_slave", "Seconds_Behind_Master")

	// Check if replication is not running:
	if ioRunning == 0 || sqlRunning == 0 {
		message += fmt.Sprintf("*IO Running:* %d\n", ioRunning)
		message += fmt.Sprintf("*SQL Running:* %d\n", sqlRunning)

		if sqlError > 0 {
			message += fmt.Sprintf("*SQL Running:* %d\n", sqlError)
		}
	} else if delay > 0 {
		message += fmt.Sprintf("*Lagging:* %d\n", delay)
		decide = true
	}

	// log.Printf("D! - Alert:MySQL:Slave - Message=%s\n", message)
	// log.Printf("D! - Alert:MySQL:Slave - Decide=%t\n", decide)

	if check == nil {
		alerts.Load().Add(
			"replication",
			"MySQL Replication Status",
			config.File.MySQL.Alerts.Replication.Duration,
			config.File.MySQL.Alerts.Replication.Warning,
			config.File.MySQL.Alerts.Replication.Critical,
			delay,
			message,
			decide,
		)
	} else {
		check.Update(delay, message)
	}
}
