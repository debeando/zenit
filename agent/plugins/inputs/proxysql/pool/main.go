package pool

import (
	"zenit/agent/plugins/inputs"
	"zenit/agent/plugins/lists/metrics"
	"zenit/config"

	"github.com/debeando/go-common/cast"
	"github.com/debeando/go-common/log"
	"github.com/debeando/go-common/mysql"
)

const SQLConnectionPool = `SELECT CASE
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

type Plugin struct{}

func (p *Plugin) Collect(name string, cnf *config.Config, mtc *metrics.Items) {
	defer func() {
		if err := recover(); err != nil {
			log.ErrorWithFields(name, log.Fields{"message": err})
		}
	}()

	for host := range cnf.Inputs.ProxySQL {
		if !cnf.Inputs.ProxySQL[host].Pool {
			return
		}

		log.DebugWithFields(name, log.Fields{
			"hostname": cnf.Inputs.ProxySQL[host].Hostname,
			"enable":   cnf.Inputs.ProxySQL[host].Enable,
			"pool":     cnf.Inputs.ProxySQL[host].Pool,
		})

		if !cnf.Inputs.ProxySQL[host].Enable {
			continue
		}

		if !cnf.Inputs.ProxySQL[host].Pool {
			continue
		}

		log.InfoWithFields(name, log.Fields{
			"hostname": cnf.Inputs.ProxySQL[host].Hostname,
		})

		m := mysql.New(cnf.Inputs.ProxySQL[host].Hostname, cnf.Inputs.ProxySQL[host].DSN)
		m.Connect()
		m.FetchAll(SQLConnectionPool, func(row map[string]string) {
			log.DebugWithFields(name, log.Fields{
				"group":    row["group"],
				"host":     row["srv_host"],
				"status":   row["status"],
				"used":     cast.StringToInt64(row["ConnUsed"]),
				"free":     cast.StringToInt64(row["ConnFree"]),
				"ok":       cast.StringToInt64(row["ConnOK"]),
				"errors":   cast.StringToInt64(row["ConnERR"]),
				"queries":  cast.StringToInt64(row["Queries"]),
				"tx":       cast.StringToInt64(row["Bytes_data_sent"]),
				"rx":       cast.StringToInt64(row["Bytes_data_recv"]),
				"latency":  cast.StringToInt64(row["Latency_us"]),
				"hostname": cnf.Inputs.ProxySQL[host].Hostname,
			})

			mtc.Add(metrics.Metric{
				Key: "proxysql_connections",
				Tags: []metrics.Tag{
					{Name: "hostname", Value: cnf.Inputs.ProxySQL[host].Hostname},
					{Name: "group", Value: row["group"]},
					{Name: "host", Value: row["srv_host"]},
				},
				Values: []metrics.Value{
					{Key: "status", Value: row["status"]},
					{Key: "used", Value: cast.StringToInt64(row["ConnUsed"])},
					{Key: "free", Value: cast.StringToInt64(row["ConnFree"])},
					{Key: "ok", Value: cast.StringToInt64(row["ConnOK"])},
					{Key: "errors", Value: cast.StringToInt64(row["ConnERR"])},
					{Key: "queries", Value: cast.StringToInt64(row["Queries"])},
					{Key: "tx", Value: cast.StringToInt64(row["Bytes_data_sent"])},
					{Key: "rx", Value: cast.StringToInt64(row["Bytes_data_recv"])},
					{Key: "latency", Value: cast.StringToInt64(row["Latency_us"])},
				},
			})
		})
	}
}

func init() {
	inputs.Add("InputProxySQLPool", func() inputs.Input { return &Plugin{} })
}
