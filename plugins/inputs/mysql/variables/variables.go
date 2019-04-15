package variables

import (
	"fmt"

	"github.com/swapbyt3s/zenit/common/log"
	"github.com/swapbyt3s/zenit/common/mysql"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/inputs"
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"
)

const query = "SHOW GLOBAL VARIABLES"

type MySQLVariables struct{}

func (l *MySQLVariables) Collect() {
	defer func() {
		if err := recover(); err != nil {
			log.Debug(fmt.Sprintf("Plugin - MySQLVariables - Panic (code %d) has been recover from somewhere.\n", err))
		}
	}()

	if !config.File.Inputs.MySQL.Variables {
		return
	}

	var a = metrics.Load()
	var m = mysql.GetInstance("mysql")
	m.Connect(config.File.Inputs.MySQL.DSN)

	rows := m.Query(query)

	for i := range rows {
		if value, ok := mysql.ParseValue(rows[i]["Value"]); ok {
			a.Add(metrics.Metric{
				Key:    "zenit_mysql_variables",
				Tags:   []metrics.Tag{{"name", rows[i]["Variable_name"]}},
				Values: value,
			})

			log.Debug(fmt.Sprintf("Plugin - InputMySQLVariables - %s=%d", rows[i]["Variable_name"], value))
		}
	}
}

func init() {
	inputs.Add("InputMySQLVariables", func() inputs.Input { return &MySQLVariables{} })
}
