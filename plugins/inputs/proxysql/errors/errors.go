package commands

import (
	"fmt"

	"github.com/swapbyt3s/zenit/common"
	"github.com/swapbyt3s/zenit/common/log"
	"github.com/swapbyt3s/zenit/common/mysql"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/inputs"
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"
)

const query = "SELECT * FROM stats_mysql_errors;"

type InputProxySQLErrors struct{}

func (l *InputProxySQLErrors) Collect() {
	defer func() {
		if err := recover(); err != nil {
			log.Debug(fmt.Sprintf("Plugin - InputProxySQLErrors - Panic (code %d) has been recover from somewhere.\n", err))
		}
	}()

	for host := range config.File.Inputs.ProxySQL {
		if !config.File.Inputs.ProxySQL[host].Errors {
			return
		}

		log.Info(fmt.Sprintf("Plugin - InputProxySQLErrors - Hostname=%s", config.File.Inputs.ProxySQL[host].Hostname))

		var a = metrics.Load()
		var p = mysql.GetInstance(config.File.Inputs.ProxySQL[host].Hostname)

		p.Connect(config.File.Inputs.ProxySQL[host].DSN)

		var r = p.Query(query)

		for _, i := range r {
			a.Add(metrics.Metric{
				Key: "proxysql_errors",
				Tags: []metrics.Tag{
					{"hostname", config.File.Inputs.ProxySQL[host].Hostname},
					{"group", i["hostgroup"]},
					{"server", i["hostname"]},
					{"port", i["port"]},
					{"username", i["username"]},
					{"schema", i["schemaname"]},
					{"errno", i["errno"]},
					{"last_error", i["last_error"]},
				},
				Values: []metrics.Value{
					{"count", common.StringToInt64(i["count_star"])},
				},
			})

			log.Debug(fmt.Sprintf("Plugin - InputProxySQLErrors - %#v", i))
		}
	}
}

func init() {
	inputs.Add("InputProxySQLErrors", func() inputs.Input { return &InputProxySQLErrors{} })
}
