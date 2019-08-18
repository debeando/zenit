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
			log.Debug(fmt.Sprintf("Plugin - InputMySQLStatus - Panic (code %d) has been recover from somewhere.\n", err))
		}
	}()

	for host := range config.File.Inputs.MySQL {
		if !config.File.Inputs.MySQL[host].Status {
			return
		}

		log.Info(fmt.Sprintf("Plugin - InputMySQLStatus - Hostname=%s", config.File.Inputs.MySQL[host].Hostname))

		var a = metrics.Load()
		var m = mysql.GetInstance("mysql")
		var v = []metrics.Value{}

		m.Connect(config.File.Inputs.MySQL[host].DSN)

		var r = m.Query(query)

		for _, i := range r {
			if value, ok := mysql.ParseValue(i["Value"]); ok {
				log.Debug(fmt.Sprintf("Plugin - InputMySQLStatus - %s=%d", i["Variable_name"], value))

				v = append(v, metrics.Value{
					Key: i["Variable_name"],
					Value: value,
				})
			}
		}

		a.Add(metrics.Metric{
			Key:    "mysql_status",
			Tags:   []metrics.Tag{
				{"hostname", config.File.Inputs.MySQL[host].Hostname},
			},
			Values: v,
		})
	}
}

func init() {
	inputs.Add("InputMySQLStatus", func() inputs.Input { return &MySQLStatus{} })
}
