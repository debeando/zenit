package queries

import (
	"regexp"
	"strings"

	"github.com/swapbyt3s/zenit/common/log"
	"github.com/swapbyt3s/zenit/common/mysql"
	"github.com/swapbyt3s/zenit/config"
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
	ReQuery = `^(?i)(SELECT|INSERT|UPDATE|DELETE)(?:.*FROM|.*INTO)?\W+([a-zA-Z0-9._]+)`
	SQL = `SELECT CASE
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

func Collect() {
	re, _ = regexp.Compile(ReQuery)

	conn, err := mysql.Connect(config.File.ProxySQL.DSN)
	defer conn.Close()
	if err != nil {
		log.Error("ProxySQL - Impossible to connect: " + err.Error())
	}

	rows, err := conn.Query(SQL)
	defer rows.Close()
	if err != nil {
		log.Error("ProxySQL - Impossible to execute query: " + err.Error())
	}

	for rows.Next() {
		var q Query

		rows.Scan(
			&q.group,
			&q.schema,
			&q.digest,
			&q.count,
			&q.sum)

		if len(q.digest) > 0 {
			table, command := Match(q.digest)

			metrics.Load().Add(metrics.Metric{
				Key: "zenit_proxysql_queries",
				Tags: []metrics.Tag{
					{"group", q.group},
					{"schema", q.schema},
					{"table", table},
					{"command", command},
				},
				Values: []metrics.Value{
					{"count", q.count},
					{"sum", q.sum},
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
