package tables

import (
	"time"

	"zenit/agent/plugins/inputs"
	"zenit/agent/plugins/lists/metrics"
	"zenit/config"

	"github.com/debeando/go-common/cast"
	"github.com/debeando/go-common/log"
	"github.com/debeando/go-common/mysql"
)

const SQLTables = `
SELECT table_schema AS 'schema',
			 table_name AS 'table',
			 COALESCE(data_length + index_length, 0) AS 'size',
			 COALESCE(table_rows, 0) AS 'rows',
			 COALESCE(auto_increment, 0) AS 'increment'
FROM information_schema.tables
WHERE table_schema NOT IN ('mysql','sys','performance_schema','information_schema')
ORDER BY table_schema, table_name;
`

type Plugin struct{}

var interval int64

func (p *Plugin) Collect(name string, cnf *config.Config, mtc *metrics.Items) {
	defer func() {
		if err := recover(); err != nil {
			log.ErrorWithFields(name, log.Fields{"message": err})
		}
	}()

	for host := range cnf.Inputs.MySQL {
		log.DebugWithFields(name, log.Fields{
			"hostname": cnf.Inputs.MySQL[host].Hostname,
			"enable":   cnf.Inputs.MySQL[host].Enable,
			"tables":   cnf.Inputs.MySQL[host].Tables.Enable,
			"interval": cnf.Inputs.MySQL[host].Tables.Interval,
		})

		if !cnf.Inputs.MySQL[host].Enable {
			continue
		}

		if !cnf.Inputs.MySQL[host].Tables.Enable {
			continue
		}

		if !IsTimeToCollect(cnf.Inputs.MySQL[host].Tables.Interval) {
			continue
		}

		log.InfoWithFields(name, log.Fields{
			"hostname": cnf.Inputs.MySQL[host].Hostname,
		})

		m := mysql.New(cnf.Inputs.MySQL[host].Hostname, cnf.Inputs.MySQL[host].DSN)
		m.Connect()
		m.FetchAll(SQLTables, func(row map[string]string) {
			log.DebugWithFields(name, log.Fields{
				"schema":    row["schema"],
				"table":     row["table"],
				"size":      row["size"],
				"rows":      row["rows"],
				"increment": row["increment"],
				"hostname":  cnf.Inputs.MySQL[host].Hostname,
			})

			mtc.Add(metrics.Metric{
				Key: "mysql_tables",
				Tags: []metrics.Tag{
					{Name: "hostname", Value: cnf.Inputs.MySQL[host].Hostname},
					{Name: "schema", Value: row["schema"]},
					{Name: "table", Value: row["table"]},
				},
				Values: []metrics.Value{
					{Key: "size", Value: cast.StringToInt64(row["size"])},
					{Key: "rows", Value: cast.StringToInt64(row["rows"])},
					{Key: "increment", Value: cast.StringToInt64(row["increment"])},
				},
			})
		})
	}
}

func init() {
	inputs.Add("InputMySQLTables", func() inputs.Input { return &Plugin{} })
}

func IsTimeToCollect(i int) bool {
	if interval == 0 || int(time.Since(time.Unix(interval, 0)).Seconds()) >= i {
		interval = int64(time.Now().Unix())
		return true
	}

	return false
}
