package pool

import (
	"fmt"

	"github.com/swapbyt3s/zenit/common/log"
	"github.com/swapbyt3s/zenit/common/mysql"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/lists/loader"
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

type InputProxySQLPool struct {}

func (l *InputProxySQLPool) Collect() {
	if ! config.File.ProxySQL.Inputs.Pool {
		return
	}

	var a = metrics.Load()
	var p = mysql.GetInstance("proxysql")
	p.Connect(config.File.ProxySQL.DSN)

	rows := p.Query(querySQLPool)

	for i := range rows {
		a.Add(metrics.Metric{
			Key: "zenit_proxysql_connection_pool",
			Tags: []metrics.Tag{
				{"group", rows[i]["group"]},
				{"host", rows[i]["srv_host"]},
			},
			Values: []metrics.Value{
				{"status", rows[i]["status"]},
				{"used", rows[i]["ConnUsed"]},
				{"free", rows[i]["ConnFree"]},
				{"ok", rows[i]["ConnOK"]},
				{"errors", rows[i]["ConnERR"]},
				{"queries", rows[i]["Queries"]},
				{"tx", rows[i]["Bytes_data_sent"]},
				{"rx", rows[i]["Bytes_data_recv"]},
				{"latency", rows[i]["Latency_us"]},
			},
		})

		log.Debug(fmt.Sprintf("Plugin - InputProxySQLPool - %#v", rows[i]))
	}
}

func init() {
	loader.Add("InputProxySQLPool", func() loader.Plugin { return &InputProxySQLPool{} })
}
