package commands

import (
	"regexp"

	"zenit/config"
	"zenit/monitor/plugins/inputs"
	"zenit/monitor/plugins/lists/metrics"

	"github.com/debeando/go-common/cast"
	"github.com/debeando/go-common/log"
	"github.com/debeando/go-common/mysql"
)

const query = "SELECT * FROM stats_mysql_errors;"

type Plugin struct{}

func (p *Plugin) Collect(name string, cnf *config.Config, mtc *metrics.Items) {
	defer func() {
		if err := recover(); err != nil {
			log.ErrorWithFields(name, log.Fields{"message": err})
		}
	}()

	for host := range cnf.Inputs.ProxySQL {
		log.DebugWithFields(name, log.Fields{
			"hostname": cnf.Inputs.ProxySQL[host].Hostname,
			"enable":   cnf.Inputs.ProxySQL[host].Enable,
			"errors":   cnf.Inputs.ProxySQL[host].Errors,
		})

		if !cnf.Inputs.ProxySQL[host].Enable {
			continue
		}

		if !cnf.Inputs.ProxySQL[host].Errors {
			continue
		}

		log.InfoWithFields(name, log.Fields{
			"hostname": cnf.Inputs.ProxySQL[host].Hostname,
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

		for _, i := range r {
			log.DebugWithFields(name, log.Fields{
				"group":      i["hostgroup"],
				"server":     i["hostname"],
				"port":       i["port"],
				"username":   i["username"],
				"schema":     i["schemaname"],
				"errno":      i["errno"],
				"last_error": parseLastError(i["last_error"]),
				"hostname":   cnf.Inputs.ProxySQL[host].Hostname,
			})

			mtc.Add(metrics.Metric{
				Key: "proxysql_errors",
				Tags: []metrics.Tag{
					{Name: "hostname", Value: cnf.Inputs.ProxySQL[host].Hostname},
					{Name: "group", Value: i["hostgroup"]},
					{Name: "server", Value: i["hostname"]},
					{Name: "port", Value: i["port"]},
					{Name: "username", Value: i["username"]},
					{Name: "schema", Value: i["schemaname"]},
					{Name: "errno", Value: i["errno"]},
					{Name: "last_error", Value: parseLastError(i["last_error"])},
				},
				Values: []metrics.Value{
					{Key: "count", Value: cast.StringToInt64(i["count_star"])},
				},
			})
		}
	}
}

func parseLastError(error string) string {
	if ok, _ := regexp.MatchString("Duplicate entry '.*' for key 'PRIMARY'", error); ok {
		error = "Duplicate entry for primary key"
	}

	return error
}

func init() {
	inputs.Add("InputProxySQLErrors", func() inputs.Input { return &Plugin{} })
}
