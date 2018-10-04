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

	var value uint64
	var ioRunning uint64
	var message string = ""
	var ok bool
	var sqlError uint64
	var sqlRunning uint64
	// var decide bool

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

	message += fmt.Sprintf("*IO Running:* %s\n", YesOrNo(ioRunning))
	message += fmt.Sprintf("*SQL Running:* %s\n", YesOrNo(sqlRunning))
	message += fmt.Sprintf("*SQL Error:* %d\n", sqlError)

	value = 2 - (ioRunning + sqlRunning)

	//log.Printf("D! - Alert:MySQL:Slave - Message=%s\n", message)
	//log.Printf("D! - Alert:MySQL:Slave - IO Running=%d\n", ioRunning)
	//log.Printf("D! - Alert:MySQL:Slave - SQL Running=%d\n", sqlRunning)
	//log.Printf("D! - Alert:MySQL:Slave - Value=%d\n", value)

	if check == nil {
		log.Printf("D! - Alert:MySQL:Slave - Adding\n")
		alerts.Load().Add(
			"replication",
			"MySQL Replication Status",
			config.File.MySQL.Alerts.Replication.Duration,
			1, // Warning
			1, // Critical
			value,
			message,
			true,
		)
	} else {
		log.Printf("D! - Alert:MySQL:Slave - Updateing\n")
		check.Update(value, message)
	}
}

func YesOrNo(v uint64) string {
	if v == 1 {
		return "Yes"
	}
	return "No"
}
