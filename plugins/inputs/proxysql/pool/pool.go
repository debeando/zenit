package pool

import (
	"fmt"

	"github.com/swapbyt3s/zenit/common"
	"github.com/swapbyt3s/zenit/common/log"
	"github.com/swapbyt3s/zenit/common/mysql"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/inputs"
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"
)

const querySQLPool = `SELECT CASE
			WHEN hostgroup IN (SELECT writer_hostgroup FROM main.mysql_replication_hostgroups) THEN 'writer'
			WHEN hostgroup IN (SELECT reader_hostgroup FROM main.mysql_replication_hostgroups) THEN 'reader'
		END AS 'group',
		srv_host,
		srv_port,
		status,
		ConnUsed,
		ConnFree,
		ConnOK,
		ConnERR,
		Queries,
		Bytes_data_sent,
		Bytes_data_recv,
		Latency_us
	FROM stats.stats_mysql_connection_pool;`

type InputProxySQLPool struct{}

func (l *InputProxySQLPool) Collect() {
	defer func() {
		if err := recover(); err != nil {
			log.Debug(fmt.Sprintf("Plugin - InputProxySQLPool - Panic (code %d) has been recover from somewhere.\n", err))
		}
	}()

	for host := range config.File.Inputs.ProxySQL {
		if !config.File.Inputs.ProxySQL[host].Pool {
			return
		}

		log.Info(fmt.Sprintf("Plugin - InputProxySQLPool - Hostname=%s", config.File.Inputs.ProxySQL[host].Hostname))

		var a = metrics.Load()
		var p = mysql.GetInstance("proxysql")

		p.Connect(config.File.Inputs.ProxySQL[host].DSN)

		var r = p.Query(querySQLPool)

		for _, i := range r {
			a.Add(metrics.Metric{
				Key: "zenit_proxysql_connections",
				Tags: []metrics.Tag{
					{"hostname", config.File.Inputs.ProxySQL[host].Hostname},
					{"group", i["group"]},
					{"host", i["srv_host"]},
				},
				Values: []metrics.Value{
					{"status", i["status"]},
					{"used", common.StringToInt64(i["ConnUsed"])},
					{"free", common.StringToInt64(i["ConnFree"])},
					{"ok", common.StringToInt64(i["ConnOK"])},
					{"errors", common.StringToInt64(i["ConnERR"])},
					{"queries", common.StringToInt64(i["Queries"])},
					{"tx", common.StringToInt64(i["Bytes_data_sent"])},
					{"rx", common.StringToInt64(i["Bytes_data_recv"])},
					{"latency", common.StringToInt64(i["Latency_us"])},
				},
			})

			log.Debug(fmt.Sprintf("Plugin - InputProxySQLPool - %#v", i))
		}
	}
}

func init() {
	inputs.Add("InputProxySQLPool", func() inputs.Input { return &InputProxySQLPool{} })
}
