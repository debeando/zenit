package tables

import (
	"fmt"

	"github.com/swapbyt3s/zenit/common"
	"github.com/swapbyt3s/zenit/common/log"
	"github.com/swapbyt3s/zenit/common/mysql"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/inputs"
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"
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
			log.Debug(fmt.Sprintf("Plugin - InputMySQLTables - Panic (code %d) has been recover from somewhere.\n", err))
		}
	}()

	for host := range config.File.Inputs.MySQL {
		if !config.File.Inputs.MySQL[host].Tables {
			return
		}

		log.Info(fmt.Sprintf("Plugin - InputMySQLTables - Hostname=%s", config.File.Inputs.MySQL[host].Hostname))


		var a = metrics.Load()
		var m = mysql.GetInstance(config.File.Inputs.MySQL[host].Hostname)

		m.Connect(config.File.Inputs.MySQL[host].DSN)

		var r = m.Query(query)

		for _, i := range r {
			log.Debug(fmt.Sprintf("Plugin - InputMySQLTables - Size %s.%s=%s", i["schema"], i["table"], i["size"]))
			log.Debug(fmt.Sprintf("Plugin - InputMySQLTables - Rows %s.%s=%s", i["schema"], i["table"], i["rows"]))
			log.Debug(fmt.Sprintf("Plugin - InputMySQLTables - Increment %s.%s=%s", i["schema"], i["table"], i["increment"]))

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
