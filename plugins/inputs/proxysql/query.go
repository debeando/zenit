package proxysql

import (
	"log"
	"regexp"
	"strings"

	"github.com/swapbyt3s/zenit/common"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/accumulator"
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
	REGEX_SQL = `^(?i)(SELECT|INSERT|UPDATE|DELETE)(?:.*FROM|.*INTO)?\W+([a-zA-Z0-9._]+)`
	QUERY_SQL = `
SELECT CASE
         WHEN hostgroup IN (SELECT writer_hostgroup FROM main.mysql_replication_hostgroups) THEN 'writer'
         WHEN hostgroup IN (SELECT reader_hostgroup FROM main.mysql_replication_hostgroups) THEN 'reader'
       END AS 'group',
       schemaname,
       digest_text,
       count_star,
       sum_time
FROM stats.stats_mysql_query_digest;
`
)

var re *regexp.Regexp

func init() {
	re, _ = regexp.Compile(REGEX_SQL)
}

func Check() bool {
	log.Printf("I! - ProxySQL - DSN: %s\n", config.ProxySQL.DSN)
	conn, err := common.MySQLConnect(config.ProxySQL.DSN)
	if err != nil {
		log.Printf("E! - ProxySQL - Impossible to connect: %s\n", err)
		return false
	}

	log.Println("I! - ProxySQL - Connected successfully.")
	conn.Close()
	return true
}

func QueryDigest() {
	conn, err := common.MySQLConnect(config.ProxySQL.DSN)
	defer conn.Close()
	if err != nil {
		panic(err)
	}

	rows, err := conn.Query(QUERY_SQL)
	defer rows.Close()
	if err != nil {
		panic(err)
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

			accumulator.Load().AddItem(accumulator.Metric{
				Key: "proxysql_queries",
				Tags: []accumulator.Tag{{"group", q.group},
					{"schema", q.schema},
					{"table", table},
					{"command", command}},
				Values: []accumulator.Value{{"count", q.count},
					{"sum", q.sum}},
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
