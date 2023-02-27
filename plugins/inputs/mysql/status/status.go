package status

import (
	"github.com/debeando/zenit/common/log"
	"github.com/debeando/zenit/common/mysql"
	"github.com/debeando/zenit/config"
	"github.com/debeando/zenit/plugins/inputs"
	"github.com/debeando/zenit/plugins/lists/metrics"
)

const query = "SHOW GLOBAL STATUS"

type MySQLStatus struct{}

func (l *MySQLStatus) Collect() {
	defer func() {
		if err := recover(); err != nil {
			log.Error("InputMySQLStatus", map[string]interface{}{"error": err})
		}
	}()

	for host := range config.File.Inputs.MySQL {
		if !config.File.Inputs.MySQL[host].Status {
			log.Debug("InputMySQLStatus", map[string]interface{}{"message": "Is not enabled."})
			return
		}

		log.Info("InputMySQLStatus", map[string]interface{}{
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
				log.Debug("InputMySQLStatus", map[string]interface{}{
					"hostname":         config.File.Inputs.MySQL[host].Hostname,
					i["Variable_name"]: value,
				})

				v = append(v, metrics.Value{Key: i["Variable_name"], Value: value})
			}
		}

		a.Add(metrics.Metric{
			Key: "mysql_status",
			Tags: []metrics.Tag{
				{Name: "hostname", Value: config.File.Inputs.MySQL[host].Hostname},
			},
			Values: v,
		})
	}
}

func init() {
	inputs.Add("InputMySQLStatus", func() inputs.Input { return &MySQLStatus{} })
}
