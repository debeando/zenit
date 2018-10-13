package errors

import (
	"fmt"

	"github.com/swapbyt3s/zenit/common"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/lists/accumulator"
	"github.com/swapbyt3s/zenit/plugins/lists/alerts"
)

func Check() {
	var metrics = accumulator.Load()
	var value = metrics.FetchOne("proxysql_connection_pool", "name", "errors")
	var errors = common.InterfaceToInt(value)

	// Build one message with details for notification:
	var message = fmt.Sprintf("*Current:* %d", errors)

	// Register new check and update last status:
	alerts.Load().Register(
		"proxysql_errors",
		"ProxySQL Connection Pool Errors",
		config.File.ProxySQL.Alerts.Errors.Duration,
		10, // Warning
		20, // Critical
		errors,
		message,
	)
}
