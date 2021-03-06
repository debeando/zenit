package variables

import (
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
			log.Error("InputMySQLVariables", map[string]interface{}{"error": err})
		}
	}()

	for host := range config.File.Inputs.MySQL {
		if !config.File.Inputs.MySQL[host].Variables {
			log.Debug("InputMySQLVariables", map[string]interface{}{"message": "Is not enabled."})
			return
		}

		log.Info("InputMySQLVariables", map[string]interface{}{
			"hostname": config.File.Inputs.MySQL[host].Hostname,
		})

		var a = metrics.Load()
		var m = mysql.GetInstance(config.File.Inputs.MySQL[host].Hostname)

		m.Connect(config.File.Inputs.MySQL[host].DSN)

		var r = m.Query(query)
		if len(r) == 0 {
			continue
		}

		var v = []metrics.Value{}

		for _, i := range r {
			if value, ok := mysql.ParseValue(i["Value"]); ok {
				log.Debug("InputMySQLVariables", map[string]interface{}{
					"hostname":         config.File.Inputs.MySQL[host].Hostname,
					i["Variable_name"]: value,
				})

				v = append(v, metrics.Value{i["Variable_name"], value})
			}
		}

		a.Add(metrics.Metric{
			Key: "mysql_variables",
			Tags: []metrics.Tag{
				{"hostname", config.File.Inputs.MySQL[host].Hostname},
			},
			Values: v,
		})
	}
}

func init() {
	inputs.Add("InputMySQLVariables", func() inputs.Input { return &MySQLVariables{} })
}
