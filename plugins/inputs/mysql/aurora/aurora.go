package aurora

import (
	"fmt"

	"github.com/swapbyt3s/zenit/common/log"
	"github.com/swapbyt3s/zenit/common/mysql"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/inputs"
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"
)

const query = "SELECT * FROM mysql.ro_replica_status"

type MySQLAurora struct{}

func (l *MySQLAurora) Collect() {
	defer func() {
		if err := recover(); err != nil {
			log.Debug(fmt.Sprintf("Plugin - InputMySQLAurora - Panic (code %d) has been recover from somewhere.\n", err))
		}
	}()

	for host := range config.File.Inputs.MySQL {
		if !config.File.Inputs.MySQL[host].Aurora {
			return
		}

		log.Info(fmt.Sprintf("Plugin - InputMySQLAurora - Hostname=%s", config.File.Inputs.MySQL[host].Hostname))

		var a = metrics.Load()
		var m = mysql.GetInstance(config.File.Inputs.MySQL[host].Hostname)
		var v = []metrics.Value{}

		m.Connect(config.File.Inputs.MySQL[host].DSN)

		var r = m.Query(query)

		for column := range r[0] {
			if value, ok := mysql.ParseValue(r[0][column]); ok {
				log.Debug(fmt.Sprintf("Plugin - InputMySQLAurora - %s=%d", column, value))

				v = append(v, metrics.Value{
					Key: column,
					Value: value,
				})
			}
		}

		a.Add(metrics.Metric{
			Key:    "aws_rds_aurora",
			Tags:   []metrics.Tag{
				{"hostname", config.File.Inputs.MySQL[host].Hostname},
			},
			Values: v,
		})
	}
}

func init() {
	inputs.Add("InputMySQLAurora", func() inputs.Input { return &MySQLAurora{} })
}
