package aurora

import (
	"fmt"
	"strconv"

	"github.com/swapbyt3s/zenit/common/log"
	"github.com/swapbyt3s/zenit/common/mysql"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/inputs"
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"
)

const query = "SELECT iops, cpu, replica_lag_in_msec FROM mysql.ro_replica_status WHERE server_id = @@aurora_server_id;"

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

		// Validar que es un aurora

		log.Info(fmt.Sprintf("Plugin - InputMySQLAurora - Hostname=%s", config.File.Inputs.MySQL[host].Hostname))

		var a = metrics.Load()
		var m = mysql.GetInstance(config.File.Inputs.MySQL[host].Hostname)
		var v = []metrics.Value{}

		m.Connect(config.File.Inputs.MySQL[host].DSN)

		var r = m.Query(query)

		if r != nil {
			for column := range r[0] {
				var val interface{}

				// buscar si hay un stringToInterface
				if isFloat(r[0][column]) {
					val, _ = strconv.ParseFloat(r[0][column], 64)
				}

				if isInt(r[0][column]) {
					if value, ok := mysql.ParseValue(r[0][column]); ok {
						val = value
					}
				}

				v = append(v, metrics.Value{
					Key: column,
					Value: val,
				})
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
}

func init() {
	inputs.Add("InputMySQLAurora", func() inputs.Input { return &MySQLAurora{} })
}

func isInt(s string) bool {
    _, err := strconv.ParseInt(s, 10, 32)
    return err == nil
}

func isFloat(s string) bool {
    _, err := strconv.ParseFloat(s, 64)
    return err == nil
}
