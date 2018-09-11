package mysql

import (
	"database/sql"
	"log"

	"github.com/swapbyt3s/zenit/common/mysql"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/accumulator"
)

const QUERY_SQL_VARIABLES = "SHOW GLOBAL VARIABLES"

func Variables() {
	conn, err := mysql.Connect(config.File.MySQL.DSN)
	defer conn.Close()
	if err != nil {
		log.Printf("E! - MySQL:Variables - Impossible to connect: %s\n", err)
	}

	rows, err := conn.Query(QUERY_SQL_VARIABLES)
	defer rows.Close()
	if err != nil {
		log.Printf("E! - MySQL:Variables - Impossible to execute query: %s\n", err)
	}

	var a = accumulator.Load()
	var k string
	var v sql.RawBytes

	for rows.Next() {
		rows.Scan(&k, &v)
		if value, ok := mysql.ParseValue(v); ok {
			a.Add(accumulator.Metric{
				Key:    "mysql_variables",
				Tags:   []accumulator.Tag{{"name", k}},
				Values: value,
			})
		}
	}
}
