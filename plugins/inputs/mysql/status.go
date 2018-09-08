package mysql

import (
	"database/sql"
	"log"

	"github.com/swapbyt3s/zenit/common/mysql"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/accumulator"
)

const QUERY_SQL_STATUS = "SHOW GLOBAL STATUS"

func Status() {
	conn, err := mysql.Connect(config.File.MySQL.DSN)
	defer conn.Close()
	if err != nil {
		log.Printf("E! - MySQL:Status - Impossible to connect: %s\n", err)
	}

	rows, err := conn.Query(QUERY_SQL_STATUS)
	defer rows.Close()
	if err != nil {
		log.Printf("E! - MySQL:Status - Impossible to execute query: %s\n", err)
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
