package lagging

import (
	"fmt"
	"log"

	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/common"
	"github.com/swapbyt3s/zenit/plugins/lists/accumulator"
	"github.com/swapbyt3s/zenit/plugins/lists/alerts"
)

func Check() {
	if ! config.File.MySQL.Inputs.Slave {
		log.Printf("W! - Require to enable MySQL Slave Status in config file.\n")
		return
	}

	var metrics = accumulator.Load()
	var value = metrics.FetchOne("mysql_slave", "name", "Seconds_Behind_Master")
	var lagging = common.InterfaceToInt(value)

	// Find own check:
	var check = alerts.Load().Exist("lagging")

	// Build one message with details for notification:
	var message = fmt.Sprintf("*Lagging:* %d\n", lagging)

	log.Printf("D! - Alert:MySQL:Slave - Message=%s\n", message)

	if check == nil {
		log.Printf("D! - Alert:MySQL:Slave - Adding\n")
		alerts.Load().Add(
			"lagging",
			"MySQL Replication Lagging",
			config.File.MySQL.Alerts.Replication.Duration,
			config.File.MySQL.Alerts.Replication.Warning,
			config.File.MySQL.Alerts.Replication.Critical,
			lagging,
			message,
			true,
		)
	} else {
		log.Printf("D! - Alert:MySQL:Slave - Updateing\n")
		check.Update(lagging, message)
	}
}
