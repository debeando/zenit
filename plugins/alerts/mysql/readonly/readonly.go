package readonly

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

type MySQLReadOnly struct {}

func (l *MySQLReadOnly) Collect() {
	if ! config.File.MySQL.Alerts.ReadOnly.Enable {
		return
	}

	if ! config.File.MySQL.Inputs.Variables {
		log.Info("Require to enable MySQL Variables in config file.")
		return
	}

	var metrics = metrics.Load()
	var value = metrics.FetchOne("zenit_mysql_variables", "name", "read_only")
	var status = common.InterfaceToUInt64(value)

	// Verify status range is valid:
	if ! (status == 0 || status == 1) {
		return
	}

	// Invert status for compatibility alert evaluation:
	status = (status ^ 1)

	// Build one message with details for notification:
	var message = fmt.Sprintf("*Current:* %s", mysql.YesOrNo(status))

	// Register new check and update last status:
	checks.Load().Register(
		"readonly",
		"MySQL Read Only",
		config.File.MySQL.Alerts.ReadOnly.Duration,
		1, // Warning
		1, // Critical
		status,
		message,
	)
}

func init() {
	alerts.Add("AlertMySQLReadOnly", func() alerts.Alert { return &MySQLReadOnly{} })
}
