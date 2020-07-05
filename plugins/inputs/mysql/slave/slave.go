package slave

import (
	"github.com/swapbyt3s/zenit/common/log"
	"github.com/swapbyt3s/zenit/common/mysql"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/inputs"
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"
)

const query = "SHOW SLAVE STATUS"

type MySQLSlave struct{}

func (l *MySQLSlave) Collect() {
	defer func() {
		if err := recover(); err != nil {
			log.Error("InputMySQLSlave", map[string]interface{}{"error": err})
		}
	}()

	for host := range config.File.Inputs.MySQL {
		if !config.File.Inputs.MySQL[host].Slave {
			return
		}

		log.Info("InputMySQLAurora", map[string]interface{}{
			"hostname": config.File.Inputs.MySQL[host].Hostname,
		})

		var a = metrics.Load()
		var m = mysql.GetInstance(config.File.Inputs.MySQL[host].Hostname)

		m.Connect(config.File.Inputs.MySQL[host].DSN)

		var r = m.Query(query)

		for column := range r[0] {
			if value, ok := mysql.ParseValue(r[0][column]); ok {
				log.Debug("InputMySQLSlave", map[string]interface{}{
					"hostname": config.File.Inputs.MySQL[host].Hostname,
					column: value,
				})

				a.Add(metrics.Metric{
					Key:    "mysql_slave",
					Tags:   []metrics.Tag{
						{"hostname", config.File.Inputs.MySQL[host].Hostname},
					},
					Values: []metrics.Value{
						{column, value},
					},
				})
			}
		}
	}
}

func init() {
	inputs.Add("InputMySQLSlave", func() inputs.Input { return &MySQLSlave{} })
}
