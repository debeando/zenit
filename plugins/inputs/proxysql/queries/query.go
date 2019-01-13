package queries

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/swapbyt3s/zenit/common"
	"github.com/swapbyt3s/zenit/common/log"
	"github.com/swapbyt3s/zenit/common/mysql"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"
	"github.com/swapbyt3s/zenit/plugins/inputs"
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
	ReQuery = `^(?i)(SELECT|INSERT|UPDATE|DELETE)(?:.*FROM|.*INTO)?\W+([a-zA-Z0-9._]+)`
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

type InputProxySQLQuery struct {}

func (l *InputProxySQLQuery) Collect() {
	if ! config.File.ProxySQL.Inputs.Queries {
		return
	}

	re, _ = regexp.Compile(ReQuery)
	var a = metrics.Load()
	var p = mysql.GetInstance("proxysql")
	p.Connect(config.File.ProxySQL.DSN)

	rows:= p.Query(querySQDigestL)

	for i := range rows {
		table, command := Match(rows[i]["digest_text"])

		if table == "unknown" || command == "unknown" {
			continue
		}

		a.Add(metrics.Metric{
			Key: "zenit_proxysql_queries",
			Tags: []metrics.Tag{
				{"group", rows[i]["group"]},
				{"schema", rows[i]["schemaname"]},
				{"table", table},
				{"command", command},
			},
			Values: []metrics.Value{
				{"count", common.StringToUInt64(rows[i]["count_star"])},
				{"sum", common.StringToUInt64(rows[i]["sum_time"])},
			},
		})

		log.Debug(
			fmt.Sprintf("Plugin - InputProxySQLQuery - (%s)%s.%s %s=%s",
				rows[i]["group"],
				rows[i]["schemaname"],
				table,
				command,
				rows[i]["count_star"],
			),
		)
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
