package pool

import (
	"fmt"

	"github.com/swapbyt3s/zenit/common"
	"github.com/swapbyt3s/zenit/common/log"
	"github.com/swapbyt3s/zenit/common/mysql"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"
	"github.com/swapbyt3s/zenit/plugins/inputs"
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

type InputProxySQLPool struct {}

func (l *InputProxySQLPool) Collect() {
  defer func () {
    if err := recover(); err != nil {
      fmt.Printf("Plugin - InputProxySQLPool - Panic (code %d) has been recover from somewhere.\n", err)
    }
  }()

	for host := range config.File.ProxySQL {
		if ! config.File.ProxySQL[host].Inputs.Pool {
			return
		}

		log.Info(fmt.Sprintf("Plugin - InputProxySQLPool - Hostname: %s", config.File.ProxySQL[host].Hostname))

		var a = metrics.Load()
		var p = mysql.GetInstance("proxysql")
		p.Connect(config.File.ProxySQL[host].DSN)

		rows := p.Query(querySQLPool)

		for i := range rows {
			a.Add(metrics.Metric{
				Key: "zenit_proxysql_connections",
				Tags: []metrics.Tag{
					{"hostname", config.File.ProxySQL[host].Hostname},
					{"group", rows[i]["group"]},
					{"host", rows[i]["srv_host"]},
				},
				Values: []metrics.Value{
					{"status", rows[i]["status"]},
					{"used", common.StringToUInt64(rows[i]["ConnUsed"])},
					{"free", common.StringToUInt64(rows[i]["ConnFree"])},
					{"ok", common.StringToUInt64(rows[i]["ConnOK"])},
					{"errors", common.StringToUInt64(rows[i]["ConnERR"])},
					{"queries", common.StringToUInt64(rows[i]["Queries"])},
					{"tx", common.StringToUInt64(rows[i]["Bytes_data_sent"])},
					{"rx", common.StringToUInt64(rows[i]["Bytes_data_recv"])},
					{"latency", common.StringToUInt64(rows[i]["Latency_us"])},
				},
			})

			log.Debug(fmt.Sprintf("Plugin - InputProxySQLPool - %#v", rows[i]))
		}
	}
}

func init() {
	inputs.Add("InputProxySQLPool", func() inputs.Input { return &InputProxySQLPool{} })
}
