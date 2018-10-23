package indexes

import (
	"strings"

	"github.com/swapbyt3s/zenit/common/log"
	"github.com/swapbyt3s/zenit/common/mysql"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"
)

// Index is a struct to save result of query.
type Index struct {
	schema      string
	table       string
	name        string
	column      string
	cardinality float64
}

const (
	querySQLIndexes = `SELECT DISTINCT
    TABLE_SCHEMA,
    TABLE_NAME,
    INDEX_NAME,
    COLUMN_NAME,
    CARDINALITY
FROM INFORMATION_SCHEMA.STATISTICS
WHERE TABLE_SCHEMA NOT IN ('mysql');`
)

func Collect() {
	conn, err := mysql.Connect(config.File.MySQL.DSN)
	defer conn.Close()
	if err != nil {
		log.Error("MySQL:Indexes - Impossible to connect: " + err.Error())
	}

	rows, err := conn.Query(querySQLIndexes)
	defer rows.Close()
	if err != nil {
		log.Error("MySQL:Indexes - Impossible to execute query: " + err.Error())
	}

	var a = metrics.Load()

	for rows.Next() {
		var i Index

		rows.Scan(
			&i.schema,
			&i.table,
			&i.name,
			&i.column,
			&i.cardinality)

		a.Add(metrics.Metric{
			Key: "zenit_mysql_indexes",
			Tags: []metrics.Tag{
				{"schema", strings.ToLower(i.schema)},
				{"table", strings.ToLower(i.table)},
				{"index", strings.ToLower(i.name)},
				{"column", strings.ToLower(i.column)},
			},
			Values: []metrics.Value{
				{"cardinality", i.cardinality},
			},
		})
	}
}
