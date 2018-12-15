package slave

import (
	"github.com/swapbyt3s/zenit/common/log"
	"github.com/swapbyt3s/zenit/common/mysql"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/lists/loader"
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"
)

const QuerySQLSlave = "SHOW SLAVE STATUS"

type MySQLSlave struct {}

func (l *MySQLSlave) Collect() {
	if ! config.File.MySQL.Inputs.Status {
		return
	}

	conn, err := mysql.Connect(config.File.MySQL.DSN)
	defer conn.Close()
	if err != nil {
		log.Error("MySQL:Slave - Impossible to connect: " + err.Error())
		return
	}

	rows, err := conn.Query(QuerySQLSlave)
	defer rows.Close()
	if err != nil {
		log.Error("MySQL:Slave - Impossible to execute query: " + err.Error())
		return
	}

	m := metrics.Load()
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
					m.Add(metrics.Metric{
						Key:    "zenit_mysql_slave",
						Tags:   []metrics.Tag{{"name", columnName}},
						Values: value,
					})
				}
			}
		}
	}
}

func init() {
	loader.Add("InputMySQLSlave", func() loader.Plugin { return &MySQLSlave{} })
}
