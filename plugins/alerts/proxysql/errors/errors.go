package errors

import (
	"fmt"

	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/lists/accumulator"
	"github.com/swapbyt3s/zenit/plugins/lists/alerts"
)

type Server struct {
	Host   string
	Group  string
	Errors int
}

func Check() {
	var metrics = accumulator.Load()
	var s Server

	for _, metric := range *metrics {
		if metric.Key == "proxysql_connection_pool" {
			for _, metricTag := range metric.Tags {
				if metricTag.Name == "host" {
					s.Host = metricTag.Value
				} else if metricTag.Name == "group" {
					s.Group = metricTag.Value
				}
			}

			for _, value := range metric.Values.([]accumulator.Value) {
				if value.Key == "errors" {
					if v, ok := value.Value.(uint); ok {
						s.Errors = int(v)
						break
					}
				}
			}

			// Build one message with details for notification:
			var message = fmt.Sprintf("*Server:* %s\n*Group:* %s\n*Error:* %d\n", s.Host, s.Group, s.Errors)

			// fmt.Printf("ProxySQL Error Message: %s\n", message)

			// Register new check and update last status:
			alerts.Load().Register(
				"proxysql_pool_errors_" + s.Host + s.Group,
				"ProxySQL Connection Pool Errors",
				config.File.ProxySQL.Alerts.Errors.Duration,
				config.File.ProxySQL.Alerts.Errors.Warning,
				config.File.ProxySQL.Alerts.Errors.Critical,
				s.Errors,
				message,
			)

		}
	}
}
