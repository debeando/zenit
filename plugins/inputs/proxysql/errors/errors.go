package commands

import (
	"regexp"

	"github.com/debeando/zenit/common"
	"github.com/debeando/zenit/common/log"
	"github.com/debeando/zenit/common/mysql"
	"github.com/debeando/zenit/config"
	"github.com/debeando/zenit/plugins/inputs"
	"github.com/debeando/zenit/plugins/lists/metrics"
)

const query = "SELECT * FROM stats_mysql_errors;"

type InputProxySQLErrors struct{}

func (l *InputProxySQLErrors) Collect() {
	defer func() {
		if err := recover(); err != nil {
			log.Error("InputProxySQLErrors", map[string]interface{}{"error": err})
		}
	}()

	for host := range config.File.Inputs.ProxySQL {
		if !config.File.Inputs.ProxySQL[host].Errors {
			return
		}

		log.Info("InputProxySQLErrors", map[string]interface{}{
			"hostname": config.File.Inputs.ProxySQL[host].Hostname,
		})

		var a = metrics.Load()
		var p = mysql.GetInstance(config.File.Inputs.ProxySQL[host].Hostname)

		p.Connect(config.File.Inputs.ProxySQL[host].DSN)

		var r = p.Query(query)

		for _, i := range r {
			log.Debug("InputProxySQLErrors", map[string]interface{}{
				"group":      i["hostgroup"],
				"server":     i["hostname"],
				"port":       i["port"],
				"username":   i["username"],
				"schema":     i["schemaname"],
				"errno":      i["errno"],
				"last_error": parseLastError(i["last_error"]),
				"hostname":   config.File.Inputs.ProxySQL[host].Hostname,
			})

			a.Add(metrics.Metric{
				Key: "proxysql_errors",
				Tags: []metrics.Tag{
					{Name: "hostname", Value: config.File.Inputs.ProxySQL[host].Hostname},
					{Name: "group", Value: i["hostgroup"]},
					{Name: "server", Value: i["hostname"]},
					{Name: "port", Value: i["port"]},
					{Name: "username", Value: i["username"]},
					{Name: "schema", Value: i["schemaname"]},
					{Name: "errno", Value: i["errno"]},
					{Name: "last_error", Value: parseLastError(i["last_error"])},
				},
				Values: []metrics.Value{
					{Key: "count", Value: common.StringToInt64(i["count_star"])},
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
	inputs.Add("InputProxySQLErrors", func() inputs.Input { return &InputProxySQLErrors{} })
}
