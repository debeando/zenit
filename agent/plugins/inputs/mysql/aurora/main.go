package aurora

import (
	"fmt"

	"zenit/agent/plugins/inputs"
	"zenit/agent/plugins/lists/metrics"
	"zenit/config"

	"github.com/debeando/go-common/log"
	"github.com/debeando/go-common/mysql"
)

type Plugin struct{}

var (
	SQLIsAurora      = "SELECT 1 FROM information_schema.TABLES WHERE (table_schema = 'mysql') AND (table_name = 'ro_replica_status')"
	SQLReplicaStatus = "SELECT * FROM mysql.ro_replica_status WHERE Server_id = '%s'"
)

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
			"aurora":   cnf.Inputs.MySQL[host].Aurora,
		})

		if !cnf.Inputs.MySQL[host].Enable {
			continue
		}

		if !cnf.Inputs.MySQL[host].Aurora {
			continue
		}

		log.InfoWithFields(name, log.Fields{
			"hostname": cnf.Inputs.MySQL[host].Hostname,
		})

		v := metrics.Values{}
		m := mysql.New(cnf.Inputs.MySQL[host].Hostname, cnf.Inputs.MySQL[host].DSN)
		m.Connect()

		if !m.FetchBool(SQLIsAurora) {
			log.WarningWithFields(name, log.Fields{
				"message":  "This server is not RDS Aurora.",
				"hostname": cnf.Inputs.MySQL[host].Hostname,
				"aurora":   cnf.Inputs.MySQL[host].Aurora,
			})
		}

		_SQLReplicaStatus := fmt.Sprintf(SQLReplicaStatus, cnf.Inputs.MySQL[host].Hostname)

		m.FetchAll(_SQLReplicaStatus, func(row map[string]string) {
			for fieldName, fieldValue := range row {
				if value, ok := mysql.ParseNumberValue(fieldValue); ok {
					log.DebugWithFields(name, log.Fields{
						"hostname": cnf.Inputs.MySQL[host].Hostname,
						fieldName: value,
					})

					v.Add(metrics.Value{Key: fieldName, Value: value})
				}
			}
		})

		mtc.Add(metrics.Metric{
			Key: "aws_aurora_rds",
			Tags: []metrics.Tag{
				{Name: "hostname", Value: cnf.Inputs.MySQL[host].Hostname},
			},
			Values: v,
		})
	}
}

func init() {
	inputs.Add("InputMySQLAurora", func() inputs.Input { return &Plugin{} })
}
