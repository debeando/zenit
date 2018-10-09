package lagging

import (
	"fmt"

	"github.com/swapbyt3s/zenit/common"
	"github.com/swapbyt3s/zenit/common/log"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/lists/accumulator"
	"github.com/swapbyt3s/zenit/plugins/lists/alerts"
)

func Check() {
	if ! config.File.MySQL.Inputs.Slave {
		log.Info("Require to enable MySQL Slave Status in config file.")
		return
	}

	var metrics = accumulator.Load()
	var value = metrics.FetchOne("mysql_slave", "name", "Seconds_Behind_Master")
	var lagging = common.InterfaceToInt(value)

	if lagging == -1 {
		return
	}

	// Find own check:
	var check = alerts.Load().Exist("lagging")

	// Build one message with details for notification:
	var message = fmt.Sprintf("*Lagging:* %d\n", lagging)

	if check == nil {
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
		check.Update(lagging, message)
	}
}
