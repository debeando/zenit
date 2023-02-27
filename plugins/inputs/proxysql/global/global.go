package global

import (
	"github.com/debeando/zenit/common/log"
	"github.com/debeando/zenit/common/mysql"
	"github.com/debeando/zenit/config"
	"github.com/debeando/zenit/plugins/inputs"
	"github.com/debeando/zenit/plugins/lists/metrics"
)

const querySQLGlobal = `SELECT Variable_Name, Variable_Value FROM stats.stats_mysql_global;`

type InputProxySQLGlobal struct{}

func (l *InputProxySQLGlobal) Collect() {
	defer func() {
		if err := recover(); err != nil {
			log.Error("InputProxySQLGlobal", map[string]interface{}{"error": err})
		}
	}()

	for host := range config.File.Inputs.ProxySQL {
		if !config.File.Inputs.ProxySQL[host].Global {
			return
		}

		log.Info("InputProxySQLGlobal", map[string]interface{}{
			"hostname": config.File.Inputs.ProxySQL[host].Hostname,
		})

		var a = metrics.Load()
		var p = mysql.GetInstance(config.File.Inputs.ProxySQL[host].Hostname)

		p.Connect(config.File.Inputs.ProxySQL[host].DSN)

		var r = p.Query(querySQLGlobal)
		if len(r) == 0 {
			continue
		}

		var v = []metrics.Value{}

		for _, i := range r {
			if value, ok := mysql.ParseValue(i["Variable_Value"]); ok {
				log.Debug("InputProxySQLGlobal", map[string]interface{}{
					i["Variable_Name"]: value,
					"hostname":         config.File.Inputs.ProxySQL[host].Hostname,
				})

				v = append(v, metrics.Value{i["Variable_Name"], value})
			}
		}

		a.Add(metrics.Metric{
			Key: "proxysql_global",
			Tags: []metrics.Tag{
				{"hostname", config.File.Inputs.ProxySQL[host].Hostname},
			},
			Values: v,
		})
	}
}

func init() {
	inputs.Add("InputProxySQLGlobal", func() inputs.Input { return &InputProxySQLGlobal{} })
}
