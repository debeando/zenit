package status

import (
	"zenit/config"
	"zenit/monitor/plugins/inputs"
	"zenit/monitor/plugins/lists/metrics"

	"github.com/debeando/go-common/log"
	"github.com/debeando/go-common/mysql"
)

const query = "SHOW GLOBAL STATUS"

type Plugin struct{}

func (p *Plugin) Collect(name string, cnf *config.Config, mtc *metrics.Items) {
	defer func() {
		if err := recover(); err != nil {
			log.ErrorWithFields(name, log.Fields{"message": err})
		}
	}()

	for host := range cnf.Inputs.MySQL {
		log.DebugWithFields(name, log.Fields{
			"hostname": cnf.Inputs.MySQL[host].Hostname,
			"enable":   cnf.Inputs.MySQL[host].Enable,
			"status":   cnf.Inputs.MySQL[host].Status,
		})

		if !cnf.Inputs.MySQL[host].Enable {
			continue
		}

		if !cnf.Inputs.MySQL[host].Status {
			continue
		}

		log.InfoWithFields(name, log.Fields{
			"hostname": cnf.Inputs.MySQL[host].Hostname,
		})

		m := mysql.New(cnf.Inputs.MySQL[host].Hostname, cnf.Inputs.MySQL[host].DSN)
		err := m.Connect()
		if err != nil {
			continue
		}

		r, _ := m.Query(query)
		if len(r) == 0 {
			continue
		}

		var v = metrics.Values{}

		for _, i := range r {
			if value, ok := mysql.ParseValue(i["Value"]); ok {
				log.DebugWithFields(name, log.Fields{
					"hostname":         cnf.Inputs.MySQL[host].Hostname,
					i["Variable_name"]: value,
				})

				v.Add(metrics.Value{Key: i["Variable_name"], Value: value})
			}
		}

		mtc.Add(metrics.Metric{
			Key: "mysql_status",
			Tags: []metrics.Tag{
				{Name: "hostname", Value: cnf.Inputs.MySQL[host].Hostname},
			},
			Values: v,
		})
	}
}

func init() {
	inputs.Add("InputMySQLStatus", func() inputs.Input { return &Plugin{} })
}
