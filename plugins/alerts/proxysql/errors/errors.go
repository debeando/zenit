package errors

import (
	"fmt"

	// "github.com/swapbyt3s/zenit/common"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/alerts"
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"
	// "github.com/swapbyt3s/zenit/plugins/lists/checks"
)

type Server struct {
	Host   string
	Group  string
	Errors uint64
}

type ProxyConnectionsErrors struct {}

func (l *ProxyConnectionsErrors) Collect() {
  defer func () {
    if err := recover(); err != nil {
      fmt.Printf("Plugin - AlertProxyConnectionsErrors - Panic (code %d) has been recover from somewhere.\n", err)
    }
  }()

	for host := range config.File.ProxySQL {
		if ! config.File.ProxySQL[host].Inputs.Commands {
			return
		}

		var m = metrics.Load()
		// var s Server

		for _, m := range *m {
			if m.Key == "zenit_proxysql_connections" {
				fmt.Printf("AlertProxyConnectionsErrors: %#v\n", m)

				// for _, metricTag := range m.Tags {
				// 	if metricTag.Name == "host" {
				// 		s.Host = metricTag.Value
				// 	} else if metricTag.Name == "group" {
				// 		s.Group = metricTag.Value
				// 	}
				// }

				// for _, value := range m.Values.([]metrics.Value) {
				// 	if value.Key == "errors" {
				// 		s.Errors = common.InterfaceToUInt64(value.Value)
				// 	}
				// }

				// // Build one message with details for notification:
				// var message = fmt.Sprintf("*Server:* %s\n*Group:* %s\n*Error:* %d\n", s.Host, s.Group, s.Errors)

				// // fmt.Printf("ProxySQL Error Message: %s\n", message)

				// // Register new check and update last status:
				// checks.Load().Register(
				// 	"proxysql_connections_errors_" + s.Host + s.Group,
				// 	"ProxySQL Connection Pool Errors",
				// 	config.File.ProxySQL[host].Alerts.Errors.Duration,
				// 	config.File.ProxySQL[host].Alerts.Errors.Warning,
				// 	config.File.ProxySQL[host].Alerts.Errors.Critical,
				// 	s.Errors,
				// 	message,
				// )
			}
		}
	}
}

func init() {
	alerts.Add("AlertProxyConnectionsErrors", func() alerts.Alert { return &ProxyConnectionsErrors{} })
}
