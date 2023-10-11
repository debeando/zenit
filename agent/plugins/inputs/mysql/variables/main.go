package variables

import (
	"zenit/agent/plugins/inputs"
	"zenit/agent/plugins/lists/metrics"
	"zenit/config"

	"github.com/debeando/go-common/log"
	"github.com/debeando/go-common/mysql"
)

const SQLShowVariables = "SHOW GLOBAL VARIABLES"

type Plugin struct{}

func (p *Plugin) Collect(name string, cnf *config.Config, mtc *metrics.Items) {
	defer func() {
		if err := recover(); err != nil {
			log.ErrorWithFields(name, log.Fields{"message": err})
		}
	}()

	for host := range cnf.Inputs.MySQL {
		log.DebugWithFields(name, log.Fields{
			"hostname":  cnf.Inputs.MySQL[host].Hostname,
			"enable":    cnf.Inputs.MySQL[host].Enable,
			"variables": cnf.Inputs.MySQL[host].Variables,
		})

		if !cnf.Inputs.MySQL[host].Enable {
			continue
		}

		if !cnf.Inputs.MySQL[host].Variables {
			continue
		}

		log.InfoWithFields(name, log.Fields{
			"hostname": cnf.Inputs.MySQL[host].Hostname,
		})

		v := metrics.Values{}
		m := mysql.New(cnf.Inputs.MySQL[host].Hostname, cnf.Inputs.MySQL[host].DSN)
		m.Connect()
		m.FetchAll(SQLShowVariables, func(row map[string]string) {
			if value, ok := mysql.ParseNumberValue(row["Value"]); ok {
				log.DebugWithFields(name, log.Fields{
					"hostname":           cnf.Inputs.MySQL[host].Hostname,
					row["Variable_name"]: value,
				})

				v.Add(metrics.Value{Key: row["Variable_name"], Value: value})
			}
		})

		mtc.Add(metrics.Metric{
			Key: "mysql_variables",
			Tags: []metrics.Tag{
				{Name: "hostname", Value: cnf.Inputs.MySQL[host].Hostname},
			},
			Values: v,
		})
	}
}

func init() {
	inputs.Add("InputMySQLVariables", func() inputs.Input { return &Plugin{} })
}
