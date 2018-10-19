package mysql

import (
	"database/sql"

	"github.com/swapbyt3s/zenit/common/log"
	"github.com/swapbyt3s/zenit/common/mysql"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"
)

const QuerySQLStatus = "SHOW GLOBAL STATUS"

func Status() {
	conn, err := mysql.Connect(config.File.MySQL.DSN)
	defer conn.Close()
	if err != nil {
		log.Error("MySQL:Status - Impossible to connect: " + err.Error())
	}

	rows, err := conn.Query(QuerySQLStatus)
	defer rows.Close()
	if err != nil {
		log.Error("MySQL:Status - Impossible to execute query: " + err.Error())
	}

	var a = metrics.Load()
	var k string
	var v sql.RawBytes

	for rows.Next() {
		err = rows.Scan(&k, &v)
		if err != nil {
			log.Error("MySQL:Slave - Error: " + err.Error())
		}

		if value, ok := mysql.ParseValue(v); ok {
			a.Add(metrics.Metric{
				Key:    "mysql_status",
				Tags:   []metrics.Tag{{"name", k}},
				Values: value,
			})
		}
	}
}
