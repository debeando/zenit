package status

import (
	"fmt"

	"github.com/swapbyt3s/zenit/common/log"
	"github.com/swapbyt3s/zenit/common/mysql"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/inputs"
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"
)

const query = "SHOW GLOBAL STATUS"

type MySQLStatus struct{}

func (l *MySQLStatus) Collect() {
	defer func() {
		if err := recover(); err != nil {
			log.Debug(fmt.Sprintf("Plugin - MySQLStatus - Panic (code %d) has been recover from somewhere.\n", err))
		}
	}()

	for host := range config.File.Inputs.MySQL {
		if !config.File.Inputs.MySQL[host].Status {
			return
		}

		log.Info(fmt.Sprintf("Plugin - MySQLStatus - Hostname: %s", config.File.Inputs.MySQL[host].Hostname))

		var a = metrics.Load()
		var m = mysql.GetInstance("mysql")
		m.Connect(config.File.Inputs.MySQL[host].DSN)

		rows := m.Query(query)

		for i := range rows {
			if value, ok := mysql.ParseValue(rows[i]["Value"]); ok {
				a.Add(metrics.Metric{
					Key:    "zenit_mysql_status",
					Tags:   []metrics.Tag{
						{"hostname", config.File.Inputs.MySQL[host].Hostname},
						{"name", rows[i]["Variable_name"]},
					},
					Values: value,
				})

				log.Debug(fmt.Sprintf("Plugin - InputMySQLStatus - %s=%d", rows[i]["Variable_name"], value))
			}
		}
	}
}

func init() {
	inputs.Add("InputMySQLStatus", func() inputs.Input { return &MySQLStatus{} })
}
