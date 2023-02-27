package queries

import (
	"regexp"
	"strings"

	"github.com/debeando/zenit/common"
	"github.com/debeando/zenit/common/log"
	"github.com/debeando/zenit/common/mysql"
	"github.com/debeando/zenit/config"
	"github.com/debeando/zenit/plugins/inputs"
	"github.com/debeando/zenit/plugins/lists/metrics"
)

type Query struct {
	group   string
	schema  string
	table   string
	command string
	digest  string
	count   uint
	sum     uint
}

const (
	ReQuery        = `^(?i)(SELECT|INSERT|UPDATE|DELETE)(?:.*FROM|.*INTO)?\W+([a-zA-Z0-9._]+)`
	querySQDigestL = `SELECT CASE
         WHEN hostgroup IN (SELECT writer_hostgroup FROM main.mysql_replication_hostgroups) THEN 'writer'
         WHEN hostgroup IN (SELECT reader_hostgroup FROM main.mysql_replication_hostgroups) THEN 'reader'
       END AS 'group',
       schemaname,
       digest_text,
       count_star,
       sum_time
FROM stats.stats_mysql_query_digest;`
)

var re *regexp.Regexp

type InputProxySQLQuery struct{}

func (l *InputProxySQLQuery) Collect() {
	defer func() {
		if err := recover(); err != nil {
			log.Error("InputProxySQLQuery", map[string]interface{}{"error": err})
		}
	}()

	for host := range config.File.Inputs.ProxySQL {
		if !config.File.Inputs.ProxySQL[host].Queries {
			return
		}

		log.Info("InputProxySQLQuery", map[string]interface{}{
			"hostname": config.File.Inputs.ProxySQL[host].Hostname,
		})

		re, _ = regexp.Compile(ReQuery)
		var a = metrics.Load()
		var p = mysql.GetInstance(config.File.Inputs.ProxySQL[host].Hostname)

		p.Connect(config.File.Inputs.ProxySQL[host].DSN)

		var r = p.Query(querySQDigestL)

		for _, i := range r {
			table, command := Match(i["digest_text"])

			if table == "unknown" || command == "unknown" {
				continue
			}

			log.Debug("InputProxySQLQuery", map[string]interface{}{
				"group":    i["group"],
				"schema":   i["schemaname"],
				"table":    table,
				"command":  command,
				"count":    common.StringToInt64(i["count_star"]),
				"sum":      common.StringToInt64(i["sum_time"]),
				"hostname": config.File.Inputs.ProxySQL[host].Hostname,
			})

			a.Add(metrics.Metric{
				Key: "proxysql_queries",
				Tags: []metrics.Tag{
					{Name: "hostname", Value: config.File.Inputs.ProxySQL[host].Hostname},
					{Name: "group", Value: i["group"]},
					{Name: "schema", Value: i["schemaname"]},
					{Name: "table", Value: table},
					{Name: "command", Value: command},
				},
				Values: []metrics.Value{
					{Key: "count", Value: common.StringToInt64(i["count_star"])},
					{Key: "sum", Value: common.StringToInt64(i["sum_time"])},
				},
			})
		}
	}
}

func Match(query string) (table string, command string) {
	groups := re.FindStringSubmatch(query)
	table = GetTable(groups)
	command = GetCommand(groups)

	return table, command
}

func IsSet(arr []string, index int) bool {
	return (len(arr) > index)
}

func GetCommand(values []string) string {
	if IsSet(values, 1) {
		return strings.ToLower(values[1])
	}
	return "unknown"
}

func GetTable(values []string) string {
	if IsSet(values, 2) {
		sql_sentence_schema_table := strings.ToLower(values[2])
		sql_sentence_objetcs := strings.Split(sql_sentence_schema_table, ".")
		return sql_sentence_objetcs[len(sql_sentence_objetcs)-1]
	}
	return "unknown"
}

func init() {
	inputs.Add("InputProxySQLQuery", func() inputs.Input { return &InputProxySQLQuery{} })
}
