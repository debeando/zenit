package status

import (
	"fmt"

	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/alerts"
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"
	"github.com/swapbyt3s/zenit/plugins/lists/checks"
)

type Server struct {
	Host       string
	Group      string
	StatusCode uint64
	StatusName string
}

type ProxyPoolStatus struct {}

func (l *ProxyPoolStatus) Collect() {
	var m = metrics.Load()
	var s Server

	for _, m := range *m {
		if m.Key == "zenit_proxysql_connections" {
			for _, metricTag := range m.Tags {
				if metricTag.Name == "host" {
					s.Host = metricTag.Value
				} else if metricTag.Name == "group" {
					s.Group = metricTag.Value
				}
			}

			for _, value := range m.Values.([]metrics.Value) {
				if value.Key == "status" {
					s.StatusCode = Status(value.Value.(string))
					s.StatusName = value.Value.(string)
					break
				}
			}

			fmt.Printf("zx %#v\n", s)

			// Build one message with details for notification:
			var message = fmt.Sprintf("*Server:* %s\n*Group:* %s\n*Status:* %s\n", s.Host, s.Group, s.StatusName)
			// fmt.Printf("ProxySQL Message: %s\n", message)

			// Register new check and update last status:
			checks.Load().Register(
				"proxysql_pool_status_" + s.Host + s.Group,
				"ProxySQL Connection Pool Status",
				config.File.ProxySQL.Alerts.Errors.Duration,
				1,
				2,
				s.StatusCode,
				message,
			)
		}
	}
}

func Status(s string) uint64 {
	switch s {
	case "ONLINE":
		return 0
	case "SHUNNED":
		return 1
	case "SHUNNED_REPLICATION_LAG":
		return 1
	case "OFFLINE_SOFT":
		return 2
	case "OFFLINE_HARD":
		return 2
	}

	return 0
}

func init() {
	alerts.Add("AlertProxyPoolStatus", func() alerts.Alert { return &ProxyPoolStatus{} })
}
