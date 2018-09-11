package mysql

import (
	"database/sql"
	"log"

	"github.com/swapbyt3s/zenit/common/mysql"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/accumulator"
)

const QUERY_SQL_SLAVE = "SHOW SLAVE STATUS"

func Slave() {
	conn, err := mysql.Connect(config.File.MySQL.DSN)
	defer conn.Close()
	if err != nil {
		log.Printf("E! - MySQL:Slave - Impossible to connect: %s\n", err)
	}

	rows, err := conn.Query(QUERY_SQL_SLAVE)
	defer rows.Close()
	if err != nil {
		log.Printf("E! - MySQL:Slave - Impossible to execute query: %s\n", err)
	}

	var a = accumulator.Load()
	var k string
	var v sql.RawBytes

	for rows.Next() {
		rows.Scan(&k, &v)
		if value, ok := mysql.ParseValue(v); ok {
			a.Add(accumulator.Metric{
				Key:    "mysql_slave",
				Tags:   []accumulator.Tag{{"name", k}},
				Values: value,
			})
		}
	}
}
