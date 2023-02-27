package tables

import (
	"github.com/debeando/zenit/common"
	"github.com/debeando/zenit/common/log"
	"github.com/debeando/zenit/common/mysql"
	"github.com/debeando/zenit/config"
	"github.com/debeando/zenit/plugins/inputs"
	"github.com/debeando/zenit/plugins/lists/metrics"
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

type MySQLTables struct{}

func (l *MySQLTables) Collect() {
	defer func() {
		if err := recover(); err != nil {
			log.Error("InputMySQLTables", map[string]interface{}{"error": err})
		}
	}()

	for host := range config.File.Inputs.MySQL {
		if !config.File.Inputs.MySQL[host].Tables {
			log.Debug("InputMySQLTables", map[string]interface{}{"message": "Is not enabled."})
			return
		}

		log.Info("InputMySQLTables", map[string]interface{}{"hostname": config.File.Inputs.MySQL[host].Hostname})

		var a = metrics.Load()
		var m = mysql.GetInstance(config.File.Inputs.MySQL[host].Hostname)

		m.Connect(config.File.Inputs.MySQL[host].DSN)

		var r = m.Query(query)

		for _, i := range r {
			log.Debug("InputMySQLTables", map[string]interface{}{
				"schema":    i["schema"],
				"table":     i["table"],
				"size":      i["size"],
				"rows":      i["rows"],
				"increment": i["increment"],
				"hostname":  config.File.Inputs.MySQL[host].Hostname,
			})

			a.Add(metrics.Metric{
				Key: "mysql_tables",
				Tags: []metrics.Tag{
					{"hostname", config.File.Inputs.MySQL[host].Hostname},
					{"schema", i["schema"]},
					{"table", i["table"]},
				},
				Values: []metrics.Value{
					{"size", common.StringToInt64(i["size"])},
					{"rows", common.StringToInt64(i["rows"])},
					{"increment", common.StringToInt64(i["increment"])},
				},
			})
		}
	}
}

func init() {
	inputs.Add("InputMySQLTables", func() inputs.Input { return &MySQLTables{} })
}
