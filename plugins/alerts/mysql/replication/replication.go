package replication

import (
	"fmt"
	"log"

	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/lists/accumulator"
	"github.com/swapbyt3s/zenit/plugins/lists/alerts"
)

func Check() {
	if ! config.File.MySQL.Inputs.Slave {
		log.Printf("W! - Require to enable MySQL Slave Status in config file.\n")
		return
	}

	var message string = ""
	var running int
	var value interface{}

	var metrics = accumulator.Load()
	var check = alerts.Load().Exist("replication")

	value = metrics.FetchOne("mysql_slave", "name", "Slave_IO_Running")
	var ioRunning = Float64ToInt(value)
	value = metrics.FetchOne("mysql_slave", "name", "Slave_SQL_Running")
	var sqlRunning = Float64ToInt(value)
	value = metrics.FetchOne("mysql_slave", "name", "Last_SQL_Errno")
	var sqlError = Float64ToInt(value)

	message += fmt.Sprintf("*IO Running:* %s\n", YesOrNo(ioRunning))
	message += fmt.Sprintf("*SQL Running:* %s\n", YesOrNo(sqlRunning))
	message += fmt.Sprintf("*SQL Error:* %d\n", sqlError)

	running = 2 - (ioRunning + sqlRunning)

	//log.Printf("D! - Alert:MySQL:Slave - Message=%s\n", message)
	//log.Printf("D! - Alert:MySQL:Slave - IO Running=%d\n", ioRunning)
	//log.Printf("D! - Alert:MySQL:Slave - SQL Running=%d\n", sqlRunning)
	//log.Printf("D! - Alert:MySQL:Slave - running=%d\n", running)

	if check == nil {
		log.Printf("D! - Alert:MySQL:Slave - Adding\n")
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
		log.Printf("D! - Alert:MySQL:Slave - Updateing\n")
		check.Update(running, message)
	}
}

func YesOrNo(v int) string {
	if v == 1 {
		return "Yes"
	}
	return "No"
}

func Float64ToInt(value interface{}) int {
	if v, ok := value.(float64); ok {
		return int(v)
	}
	return -1
}
