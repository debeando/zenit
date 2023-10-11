package replica

import (
	"zenit/agent/plugins/inputs"
	"zenit/agent/plugins/lists/metrics"
	"zenit/config"

	"github.com/debeando/go-common/log"
	"github.com/debeando/go-common/mysql"
)

// TODO:
//   - This command in future version is deprecated.
//     Use: SHOW REPLICA STATUS
//   - You need detect version.
const SQLShowReplicaStatus = "SHOW SLAVE STATUS"

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
			"replica":  cnf.Inputs.MySQL[host].Replica,
		})

		if !cnf.Inputs.MySQL[host].Enable {
			continue
		}

		if !cnf.Inputs.MySQL[host].Replica {
			continue
		}

		log.InfoWithFields(name, log.Fields{
			"hostname": cnf.Inputs.MySQL[host].Hostname,
		})

		v := metrics.Values{}
		m := mysql.New(cnf.Inputs.MySQL[host].Hostname, cnf.Inputs.MySQL[host].DSN)
		m.Connect()
		m.FetchAll(SQLShowReplicaStatus, func(row map[string]string) {
			for column := range row {
				if value, ok := mysql.ParseValue(row[column]); ok {
					log.DebugWithFields(name, log.Fields{
						"hostname": cnf.Inputs.MySQL[host].Hostname,
						column:     value,
					})

					v.Add(metrics.Value{Key: column, Value: value})
				}
			}
		})

		mtc.Add(metrics.Metric{
			Key: "mysql_slave",
			Tags: []metrics.Tag{
				{Name: "hostname", Value: cnf.Inputs.MySQL[host].Hostname},
			},
			Values: v,
		})
	}
}

func init() {
	inputs.Add("InputMySQLReplica", func() inputs.Input { return &Plugin{} })
}
