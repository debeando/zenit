package mysql

import (
	"database/sql"

	"github.com/swapbyt3s/zenit/common/mysql"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/accumulator"
)

const QUERY_SQL_STATUS = "SHOW GLOBAL STATUS"

func Status() {
	conn, err := mysql.Connect(config.MySQL.DSN)
	defer conn.Close()
	if err != nil {
		panic(err)
	}

	rows, err := conn.Query(QUERY_SQL_STATUS)
	defer rows.Close()
	if err != nil {
		panic(err)
	}

	var a = accumulator.Load()
	var k string
	var v sql.RawBytes

	for rows.Next() {
		rows.Scan(&k, &v)
		if value, ok := mysql.ParseValue(v); ok {
			a.AddItem(accumulator.Metric{
				Key:    "mysql_status",
				Tags:   []accumulator.Tag{{"name", k}},
				Values: value,
			})
		}
	}
}
