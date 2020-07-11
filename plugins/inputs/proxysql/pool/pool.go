package pool

import (
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
		(SELECT max_connections FROM mysql_servers WHERE hostname = srv_host AND hostgroup_id = hostgroup) AS ConnMax,
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
			log.Error("InputProxySQLPool", map[string]interface{}{"error": err})
		}
	}()

	for host := range config.File.Inputs.ProxySQL {
		if !config.File.Inputs.ProxySQL[host].Pool {
			return
		}

		log.Info("InputProxySQLPool", map[string]interface{}{
			"hostname": config.File.Inputs.ProxySQL[host].Hostname,
		})

		var a = metrics.Load()
		var p = mysql.GetInstance(config.File.Inputs.ProxySQL[host].Hostname)

		p.Connect(config.File.Inputs.ProxySQL[host].DSN)

		var r = p.Query(querySQLPool)

		for _, i := range r {
			log.Debug("InputProxySQLPool", map[string]interface{}{
				"group":    i["group"],
				"host":     i["srv_host"],
				"status":   i["status"],
				"used":     common.StringToInt64(i["ConnUsed"]),
				"free":     common.StringToInt64(i["ConnFree"]),
				"ok":       common.StringToInt64(i["ConnOK"]),
				"errors":   common.StringToInt64(i["ConnERR"]),
				"queries":  common.StringToInt64(i["Queries"]),
				"tx":       common.StringToInt64(i["Bytes_data_sent"]),
				"rx":       common.StringToInt64(i["Bytes_data_recv"]),
				"latency":  common.StringToInt64(i["Latency_us"]),
				"hostname": config.File.Inputs.ProxySQL[host].Hostname,
			})

			a.Add(metrics.Metric{
				Key: "proxysql_connections",
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
		}
	}
}

func init() {
	inputs.Add("InputProxySQLPool", func() inputs.Input { return &InputProxySQLPool{} })
}
