package readonly

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
	if ! config.File.MySQL.Inputs.Variables {
		log.Info("Require to enable MySQL Variables in config file.")
		return
	}

	var metrics = accumulator.Load()
	var value = metrics.FetchOne("mysql_variables", "name", "read_only")
	var status = common.InterfaceToInt(value)

	// Verify status range is valid:
	if ! (status == 0 || status == 1) {
		return
	}

	// Invert status for compatibility alert evaluation:
	status = (status ^ 1)

	// Build one message with details for notification:
	var message = fmt.Sprintf("*Current:* %s", mysql.YesOrNo(status ^ 1))

	// Register new check and update last status:
	alerts.Load().Register(
		"readonly",
		"MySQL Read Only",
		config.File.MySQL.Alerts.ReadOnly.Duration,
		1, // Warning
		1, // Critical
		status,
		message,
	)
}
