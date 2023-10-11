package commands

import (
	"regexp"

	"zenit/agent/plugins/inputs"
	"zenit/agent/plugins/lists/metrics"
	"zenit/config"

	"github.com/debeando/go-common/cast"
	"github.com/debeando/go-common/log"
	"github.com/debeando/go-common/mysql"
)

const SQLErrors = "SELECT * FROM stats_mysql_errors;"

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
		m.Connect()
		m.FetchAll(SQLErrors, func(row map[string]string) {
			log.DebugWithFields(name, log.Fields{
				"group":      row["hostgroup"],
				"server":     row["hostname"],
				"port":       row["port"],
				"username":   row["username"],
				"schema":     row["schemaname"],
				"errno":      row["errno"],
				"last_error": parseLastError(row["last_error"]),
				"hostname":   cnf.Inputs.ProxySQL[host].Hostname,
			})

			mtc.Add(metrics.Metric{
				Key: "proxysql_errors",
				Tags: []metrics.Tag{
					{Name: "hostname", Value: cnf.Inputs.ProxySQL[host].Hostname},
					{Name: "group", Value: row["hostgroup"]},
					{Name: "server", Value: row["hostname"]},
					{Name: "port", Value: row["port"]},
					{Name: "username", Value: row["username"]},
					{Name: "schema", Value: row["schemaname"]},
					{Name: "errno", Value: row["errno"]},
					{Name: "last_error", Value: parseLastError(row["last_error"])},
				},
				Values: []metrics.Value{
					{Key: "count", Value: cast.StringToInt64(row["count_star"])},
				},
			})
		})
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
