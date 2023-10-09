package tables

import (
	"zenit/config"
	"zenit/monitor/plugins/inputs"
	"zenit/monitor/plugins/lists/metrics"

	"github.com/debeando/go-common/cast"
	"github.com/debeando/go-common/log"
	"github.com/debeando/go-common/mysql"
)

const query = `
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

func (p *Plugin) Collect(name string, cnf *config.Config, mtc *metrics.Items) {
	defer func() {
		if err := recover(); err != nil {
			log.ErrorWithFields(name, log.Fields{"error": err})
		}
	}()

	for host := range cnf.Inputs.MySQL {
		log.DebugWithFields(name, log.Fields{
			"hostname": cnf.Inputs.MySQL[host].Hostname,
			"enable":   cnf.Inputs.MySQL[host].Enable,
			"tables":   cnf.Inputs.MySQL[host].Tables,
		})

		if !cnf.Inputs.MySQL[host].Enable {
			continue
		}

		if !cnf.Inputs.MySQL[host].Tables {
			continue
		}

		log.InfoWithFields(name, log.Fields{
			"hostname": cnf.Inputs.MySQL[host].Hostname,
		})

		m := mysql.New(cnf.Inputs.MySQL[host].Hostname, cnf.Inputs.MySQL[host].DSN)
		m.Connect()

		r, _ := m.Query(query)
		if len(r) == 0 {
			continue
		}

		for _, i := range r {
			log.DebugWithFields(name, log.Fields{
				"schema":    i["schema"],
				"table":     i["table"],
				"size":      i["size"],
				"rows":      i["rows"],
				"increment": i["increment"],
				"hostname":  cnf.Inputs.MySQL[host].Hostname,
			})

			mtc.Add(metrics.Metric{
				Key: "mysql_tables",
				Tags: []metrics.Tag{
					{Name: "hostname", Value: cnf.Inputs.MySQL[host].Hostname},
					{Name: "schema", Value: i["schema"]},
					{Name: "table", Value: i["table"]},
				},
				Values: []metrics.Value{
					{Key: "size", Value: cast.StringToInt64(i["size"])},
					{Key: "rows", Value: cast.StringToInt64(i["rows"])},
					{Key: "increment", Value: cast.StringToInt64(i["increment"])},
				},
			})
		}
	}
}

func init() {
	inputs.Add("InputMySQLTables", func() inputs.Input { return &Plugin{} })
}
