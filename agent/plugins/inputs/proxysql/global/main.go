package global

import (
	"zenit/config"
	"zenit/agent/plugins/inputs"
	"zenit/agent/plugins/lists/metrics"

	"github.com/debeando/go-common/log"
	"github.com/debeando/go-common/mysql"
)

const query = `SELECT Variable_Name, Variable_Value FROM stats.stats_mysql_global;`

type Plugin struct{}

func (p *Plugin) Collect(name string, cnf *config.Config, mtc *metrics.Items) {
	defer func() {
		if err := recover(); err != nil {
			log.ErrorWithFields(name, log.Fields{"message": err})
		}
	}()

	for host := range cnf.Inputs.ProxySQL {
		if !cnf.Inputs.ProxySQL[host].Global {
			return
		}

		log.DebugWithFields(name, log.Fields{
			"hostname": cnf.Inputs.ProxySQL[host].Hostname,
			"enable":   cnf.Inputs.ProxySQL[host].Enable,
			"global":   cnf.Inputs.ProxySQL[host].Global,
		})

		if !cnf.Inputs.ProxySQL[host].Enable {
			continue
		}

		if !cnf.Inputs.ProxySQL[host].Global {
			continue
		}

		log.InfoWithFields(name, log.Fields{
			"hostname": cnf.Inputs.ProxySQL[host].Hostname,
		})

		var v = metrics.Values{}

		m := mysql.New(cnf.Inputs.MySQL[host].Hostname, cnf.Inputs.MySQL[host].DSN)
		err := m.Connect()
		if err != nil {
			continue
		}

		r, _ := m.Query(query)
		if len(r) == 0 {
			continue
		}

		for _, i := range r {
			if value, ok := mysql.ParseValue(i["Variable_Value"]); ok {
				log.DebugWithFields(name, log.Fields{
					i["Variable_Name"]: value,
					"hostname":         cnf.Inputs.ProxySQL[host].Hostname,
				})

				v.Add(metrics.Value{Key: i["Variable_Name"], Value: value})
			}
		}

		mtc.Add(metrics.Metric{
			Key: "proxysql_global",
			Tags: []metrics.Tag{
				{Name: "hostname", Value: cnf.Inputs.ProxySQL[host].Hostname},
			},
			Values: v,
		})
	}
}

func init() {
	inputs.Add("InputProxySQLGlobal", func() inputs.Input { return &Plugin{} })
}
