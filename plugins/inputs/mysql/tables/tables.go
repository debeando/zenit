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
WHERE table_schema NOT IN ('mysql','sys','performance_schema','information_schema','percona')
ORDER BY table_schema, table_name;
`

type MySQLTables struct{}

func (l *MySQLTables) Collect() {
	defer func() {
		if err := recover(); err != nil {
			log.Debug(fmt.Sprintf("Plugin - MySQLTables - Panic (code %d) has been recover from somewhere.\n", err))
		}
	}()

	if !config.File.MySQL.Inputs.Tables {
		return
	}

	var a = metrics.Load()
	var m = mysql.GetInstance("mysql")
	m.Connect(config.File.MySQL.DSN)

	rows := m.Query(query)

	for i := range rows {
		a.Add(metrics.Metric{
			Key: "zenit_mysql_tables",
			Tags: []metrics.Tag{
				{"schema", rows[i]["schema"]},
				{"table", rows[i]["table"]},
			},
			Values: []metrics.Value{
				{"size", common.StringToUInt64(rows[i]["size"])},
				{"rows", common.StringToUInt64(rows[i]["rows"])},
				{"increment", common.StringToUInt64(rows[i]["increment"])}},
		})

		log.Debug(fmt.Sprintf("Plugin - InputMySQLTables - Size %s.%s=%s", rows[i]["schema"], rows[i]["table"], rows[i]["size"]))
		log.Debug(fmt.Sprintf("Plugin - InputMySQLTables - Rows %s.%s=%s", rows[i]["schema"], rows[i]["table"], rows[i]["rows"]))
		log.Debug(fmt.Sprintf("Plugin - InputMySQLTables - Increment %s.%s=%s", rows[i]["schema"], rows[i]["table"], rows[i]["increment"]))
	}
}

func init() {
	inputs.Add("InputMySQLTables", func() inputs.Input { return &MySQLTables{} })
}
