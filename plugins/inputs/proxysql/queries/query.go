package queries

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/swapbyt3s/zenit/common"
	"github.com/swapbyt3s/zenit/common/log"
	"github.com/swapbyt3s/zenit/common/mysql"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/inputs"
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"
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
			log.Debug(fmt.Sprintf("Plugin - InputProxySQLQuery - Panic (code %d) has been recover from somewhere.\n", err))
		}
	}()

	for host := range config.File.Inputs.ProxySQL {
		if !config.File.Inputs.ProxySQL[host].Queries {
			return
		}

		log.Info(fmt.Sprintf("Plugin - InputProxySQLQuery - Hostname: %s", config.File.Inputs.ProxySQL[host].Hostname))

		re, _ = regexp.Compile(ReQuery)
		var a = metrics.Load()
		var p = mysql.GetInstance("proxysql")

		p.Connect(config.File.Inputs.ProxySQL[host].DSN)

		var r = p.Query(querySQDigestL)

		for _, i := range r {
			table, command := Match(i["digest_text"])

			if table == "unknown" || command == "unknown" {
				continue
			}

			a.Add(metrics.Metric{
				Key: "zenit_proxysql_queries",
				Tags: []metrics.Tag{
					{"group", i["group"]},
					{"schema", i["schemaname"]},
					{"table", table},
					{"command", command},
				},
				Values: []metrics.Value{
					{"count", common.StringToInt64(i["count_star"])},
					{"sum", common.StringToInt64(i["sum_time"])},
				},
			})

			log.Debug(
				fmt.Sprintf("Plugin - InputProxySQLQuery - (%s)%s.%s %s=%s",
					i["group"],
					i["schemaname"],
					table,
					command,
					i["count_star"],
				),
			)
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
