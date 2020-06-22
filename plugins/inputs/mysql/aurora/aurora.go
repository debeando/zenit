package aurora

import (
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
			log.Error("InputMySQLAurora", map[string]interface{}{"error": err})
		}
	}()

	for host := range config.File.Inputs.MySQL {
		if !config.File.Inputs.MySQL[host].Aurora {
			return
		}

		log.Info("InputMySQLAurora", map[string]interface{}{"hostname": config.File.Inputs.MySQL[host].Hostname})

		var a = metrics.Load()
		var m = mysql.GetInstance(config.File.Inputs.MySQL[host].Hostname)
		var v = []metrics.Value{}

		m.Connect(config.File.Inputs.MySQL[host].DSN)

		var c = m.Query("SELECT 1 FROM information_schema.TABLES WHERE (TABLE_SCHEMA = 'mysql') AND (TABLE_NAME = 'ro_replica_status')")
		if len(c) == 0 {
			continue
		}

		var r = m.Query(query)

		if len(r) == 0 {
			continue
		}

		for column := range r[0] {
			if value, ok := mysql.ParseValue(r[0][column]); ok {
				log.Debug("InputMySQLAurora", map[string]interface{}{"attribute": column, "value": value})

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
