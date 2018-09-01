package mysql

import (
	"log"

	"github.com/swapbyt3s/zenit/common/mysql"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/accumulator"
)

type Table struct {
	schema    string
	table     string
	size      float64
	rows      float64
	increment float64
}

const QUERY_SQL_TABLES = `
SELECT table_schema AS 'schema',
       table_name AS 'table',
       data_length + index_length AS 'size',
       table_rows AS 'rows',
       auto_increment AS 'increment'
FROM information_schema.tables
WHERE table_schema NOT IN ('mysql','sys','performance_schema','information_schema','percona')
ORDER BY table_schema, table_name;
`

func Tables() {
	conn, err := mysql.Connect(config.MySQL.DSN)
	defer conn.Close()
	if err != nil {
		log.Printf("E! - MySQL:Tables - Impossible to connect: %s\n", err)
	}

	rows, err := conn.Query(QUERY_SQL_TABLES)
	defer rows.Close()
	if err != nil {
		log.Printf("E! - MySQL:Tables - Impossible to execute query: %s\n", err)
	}

	var a = accumulator.Load()

	for rows.Next() {
		var t Table

		rows.Scan(
			&t.schema,
			&t.table,
			&t.size,
			&t.rows,
			&t.increment)

		a.AddItem(accumulator.Metric{
			Key: "mysql_stats_tables",
			Tags: []accumulator.Tag{{"schema", t.schema},
				{"table", t.table}},
			Values: []accumulator.Value{{"size", t.size},
				{"rows", t.rows},
				{"increment", t.increment}},
		})
	}
}
