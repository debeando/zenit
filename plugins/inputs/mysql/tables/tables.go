package tables

import (
	// "fmt"

	// "github.com/swapbyt3s/zenit/common/log"
	"github.com/swapbyt3s/zenit/common/mysql"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/lists/loader"
	// "github.com/swapbyt3s/zenit/plugins/lists/metrics"
)

type Table struct {
	schema    string
	table     string
	size      float64
	rows      float64
	increment float64
}

const querySQLTable = `
SELECT table_schema AS 'schema',
       table_name AS 'table',
       data_length + index_length AS 'size',
       table_rows AS 'rows',
       auto_increment AS 'increment'
FROM information_schema.tables
WHERE table_schema NOT IN ('mysql','sys','performance_schema','information_schema','percona')
ORDER BY table_schema, table_name;
`

type MySQLTables struct {}

func (l *MySQLTables) Collect() {
	if ! config.File.MySQL.Inputs.Tables {
		return
	}

	var m = mysql.GetInstance("mysql")
	m.Connect(config.File.MySQL.DSN)

	m.Query(querySQLTable)
	// fmt.Printf("%#v", rows)

//	var a = metrics.Load()
//
//	for rows.Next() {
//		var t Table
//
//		rows.Scan(
//			&t.schema,
//			&t.table,
//			&t.size,
//			&t.rows,
//			&t.increment)
//
//		a.Add(metrics.Metric{
//			Key: "zenit_mysql_stats_tables",
//			Tags: []metrics.Tag{
//				{"schema", t.schema},
//				{"table", t.table}},
//			Values: []metrics.Value{
//				{"size", uint(t.size)},
//				{"rows", uint(t.rows)},
//				{"increment", uint(t.increment)}},
//		})
//	}
}

func init() {
	loader.Add("InputMySQLTables", func() loader.Plugin { return &MySQLTables{} })
}
