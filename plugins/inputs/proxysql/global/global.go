package global

import (
	"github.com/swapbyt3s/zenit/common/log"
	"github.com/swapbyt3s/zenit/common/mysql"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/inputs"
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"
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

		for _, i := range r {
			if value, ok := mysql.ParseValue(i["Variable_Value"]); ok {
				log.Debug("InputProxySQLGlobal", map[string]interface{}{
					"attribute": i["Variable_name"],
					"value": value,
					"hostname": config.File.Inputs.ProxySQL[host].Hostname,
				})

				a.Add(metrics.Metric{
					Key:    "proxysql_global",
					Tags:   []metrics.Tag{
						{"hostname", config.File.Inputs.ProxySQL[host].Hostname},
					},
					Values: []metrics.Value{
						{i["Variable_Name"], value},
					},
				})
			}
		}
	}
}

func init() {
	inputs.Add("InputProxySQLGlobal", func() inputs.Input { return &InputProxySQLGlobal{} })
}
