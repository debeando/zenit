package queries

import (
	"regexp"
	"strings"

	"zenit/agent/plugins/inputs"
	"zenit/agent/plugins/lists/metrics"
	"zenit/config"

	"github.com/debeando/go-common/cast"
	"github.com/debeando/go-common/log"
	"github.com/debeando/go-common/mysql"
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
	ReQuery    = `^(?i)(SELECT|INSERT|UPDATE|DELETE)(?:.*FROM|.*INTO)?\W+([a-zA-Z0-9._]+)`
	SQLQueries = `SELECT CASE
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

type Plugin struct{}

func (p *Plugin) Collect(name string, cnf *config.Config, mtc *metrics.Items) {
	defer func() {
		if err := recover(); err != nil {
			log.ErrorWithFields(name, log.Fields{"message": err})
		}
	}()

	for host := range cnf.Inputs.ProxySQL {
		if !cnf.Inputs.ProxySQL[host].Queries {
			return
		}

		log.DebugWithFields(name, log.Fields{
			"hostname": cnf.Inputs.ProxySQL[host].Hostname,
			"enable":   cnf.Inputs.ProxySQL[host].Enable,
			"queries":  cnf.Inputs.ProxySQL[host].Queries,
		})

		if !cnf.Inputs.ProxySQL[host].Enable {
			continue
		}

		if !cnf.Inputs.ProxySQL[host].Queries {
			continue
		}

		log.InfoWithFields(name, log.Fields{
			"hostname": cnf.Inputs.ProxySQL[host].Hostname,
		})

		re, _ = regexp.Compile(ReQuery)

		m := mysql.New(cnf.Inputs.ProxySQL[host].Hostname, cnf.Inputs.ProxySQL[host].DSN)
		m.Connect()
		m.FetchAll(SQLQueries, func(row map[string]string) {
			table, command := Match(row["digest_text"])

			if table == "unknown" || command == "unknown" {
				return
			}

			log.DebugWithFields(name, log.Fields{
				"group":    row["group"],
				"schema":   row["schemaname"],
				"table":    table,
				"command":  command,
				"count":    cast.StringToInt64(row["count_star"]),
				"sum":      cast.StringToInt64(row["sum_time"]),
				"hostname": cnf.Inputs.ProxySQL[host].Hostname,
			})

			mtc.Add(metrics.Metric{
				Key: "proxysql_queries",
				Tags: []metrics.Tag{
					{Name: "hostname", Value: cnf.Inputs.ProxySQL[host].Hostname},
					{Name: "group", Value: row["group"]},
					{Name: "schema", Value: row["schemaname"]},
					{Name: "table", Value: table},
					{Name: "command", Value: command},
				},
				Values: []metrics.Value{
					{Key: "count", Value: cast.StringToInt64(row["count_star"])},
					{Key: "sum", Value: cast.StringToInt64(row["sum_time"])},
				},
			})
		})
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
	inputs.Add("InputProxySQLQuery", func() inputs.Input { return &Plugin{} })
}
