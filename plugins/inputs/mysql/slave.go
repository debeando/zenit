package mysql

import (
	"github.com/swapbyt3s/zenit/common/log"
	"github.com/swapbyt3s/zenit/common/mysql"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/lists/accumulator"
)

const QuerySQLSlave = "SHOW SLAVE STATUS"

func Slave() {
	conn, err := mysql.Connect(config.File.MySQL.DSN)
	defer conn.Close()
	if err != nil {
		log.Error("MySQL:Slave - Impossible to connect: " + err.Error())
	}

	rows, err := conn.Query(QuerySQLSlave)
	defer rows.Close()
	if err != nil {
		log.Error("MySQL:Slave - Impossible to execute query: " + err.Error())
	}

	metrics := accumulator.Load()
	columns, _ := rows.Columns()
	count := len(columns)
	status := make([]interface{}, count)
	values := make([]interface{}, count)

	for rows.Next() {
		for columnIndex, _ := range columns {
			values[columnIndex] = &status[columnIndex]
		}

		err = rows.Scan(values...)
		if err != nil {
			log.Error("MySQL:Slave - Error: " + err.Error())
		}

		for columnIndex, columnName := range columns {
			if state, ok := status[columnIndex].([]byte); ok {
				if value, ok := mysql.ParseValue(state); ok {
					metrics.Add(accumulator.Metric{
						Key:    "mysql_slave",
						Tags:   []accumulator.Tag{{"name", columnName}},
						Values: value,
					})
				}
			}
		}
	}
}
