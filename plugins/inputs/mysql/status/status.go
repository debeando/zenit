package status

import (
	"fmt"

	"github.com/swapbyt3s/zenit/common/log"
	"github.com/swapbyt3s/zenit/common/mysql"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"
	"github.com/swapbyt3s/zenit/plugins/inputs"
)

const query = "SHOW GLOBAL STATUS"

type MySQLStatus struct {}

func (l *MySQLStatus) Collect() {
	if ! config.File.MySQL.Inputs.Status {
		return
	}

	var a = metrics.Load()
	var m = mysql.GetInstance("mysql")
	m.Connect(config.File.MySQL.DSN)

	rows := m.Query(query)

	for i := range rows {
		if value, ok := mysql.ParseValue(rows[i]["Value"]); ok {
			a.Add(metrics.Metric{
				Key:    "zenit_mysql_status",
				Tags:   []metrics.Tag{{"name", rows[i]["Variable_name"]}},
				Values: value,
			})

			log.Debug(fmt.Sprintf("Plugin - InputMySQLStatus - %s=%d", rows[i]["Variable_name"], value))
		}
	}
}

func init() {
	inputs.Add("InputMySQLStatus", func() inputs.Input { return &MySQLStatus{} })
}
