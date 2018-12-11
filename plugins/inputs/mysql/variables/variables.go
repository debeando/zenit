package variables

import (
	"database/sql"

	"github.com/swapbyt3s/zenit/common/log"
	"github.com/swapbyt3s/zenit/common/mysql"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/lists/loader"
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"
)

const QuerySQLVariables = "SHOW GLOBAL VARIABLES"

type MySQLVariables struct {}

func (l *MySQLVariables) Collect() {
	if ! config.File.MySQL.Inputs.Variables {
		return
	}

	conn, err := mysql.Connect(config.File.MySQL.DSN)
	defer conn.Close()
	if err != nil {
		log.Error("MySQL:Variables - Impossible to connect: " + err.Error())
	}

	rows, err := conn.Query(QuerySQLVariables)
	defer rows.Close()
	if err != nil {
		log.Error("MySQL:Variables - Impossible to execute query: " + err.Error())
	}

	var a = metrics.Load()
	var k string
	var v sql.RawBytes

	for rows.Next() {
		rows.Scan(&k, &v)
		if value, ok := mysql.ParseValue(v); ok {
			a.Add(metrics.Metric{
				Key:    "zenit_mysql_variables",
				Tags:   []metrics.Tag{{"name", k}},
				Values: value,
			})
		}
	}
}

func init() {
	loader.Add("MySQLVariables", func() loader.Plugin { return &MySQLVariables{} })
}
