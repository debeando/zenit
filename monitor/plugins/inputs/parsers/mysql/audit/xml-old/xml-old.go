// This parse OLD Style
// https://dev.mysql.com/doc/refman/5.5/en/audit-log-file-formats.html
// TODO: Move this package to inputs/parsers/mysqlauditlog

package xmlold

import (
	"zenit/config"
	"zenit/monitor/plugins/outputs/clickhouse"

	"github.com/debeando/go-common/cast"
	"github.com/debeando/go-common/maps"
	"github.com/debeando/go-common/mysql"
	"github.com/debeando/go-common/mysql/sql"
	"github.com/debeando/go-common/strings"
	"github.com/debeando/go-common/tail"
)

func Collect() {
	cnf := config.GetInstance()
	channel_tail := make(chan string)
	channel_parser := make(chan map[string]string)
	channel_data := make(chan map[string]string)

	event := &clickhouse.Event{
		Type:    "AuditLog",
		Schema:  "zenit",
		Table:   "mysql_audit_log",
		Size:    cnf.Parser.MySQL.AuditLog.BufferSize,
		Timeout: cnf.Parser.MySQL.AuditLog.BufferTimeOut,
		Wildcard: map[string]string{
			"_time":          "'%s'",
			"command_class":  "'%s'",
			"connection_id":  "%s",
			"db":             "'%s'",
			"host":           "'%s'",
			"host_ip":        "IPv4StringToNum('%s')",
			"host_name":      "'%s'",
			"ip":             "'%s'",
			"name":           "'%s'",
			"os_login":       "'%s'",
			"os_user":        "'%s'",
			"priv_user":      "'%s'",
			"proxy_user":     "'%s'",
			"record":         "'%s'",
			"sqltext":        "'%s'",
			"sqltext_digest": "'%s'",
			"status":         "%s",
			"user":           "'%s'",
		},
		Values: []map[string]string{},
	}

	go tail.Tail(cnf.Parser.MySQL.AuditLog.LogPath, channel_tail)
	go Parser(channel_tail, channel_parser)
	go clickhouse.Send(event, channel_data)

	go func() {
		for channel_event := range channel_parser {
			channel_data <- channel_event
		}
	}()
}

func Parser(tail <-chan string, parser chan<- map[string]string) {
	cnf := config.GetInstance()
	var buffer []string

	go func() {
		defer close(parser)

		for line := range tail {
			if line == "<AUDIT_RECORD" && len(buffer) > 0 {
				result := map[string]string{
					"_time":          "",
					"command_class":  "",
					"connection_id":  "",
					"db":             "",
					"host":           "",
					"host_ip":        "",
					"host_name":      "",
					"ip":             "",
					"name":           "",
					"os_login":       "",
					"os_user":        "",
					"priv_user":      "",
					"proxy_user":     "",
					"record":         "",
					"sqltext":        "",
					"sqltext_digest": "",
					"status":         "",
					"user":           "",
				}

				for _, kv := range buffer {
					key, value := strings.SplitKeyAndValue(&kv)
					result[key] = strings.Trim(&value)
				}

				buffer = buffer[:0]

				if maps.In("user", result) {
					result["user"] = mysql.ClearUser(result["user"])
				}

				if maps.In("sqltext", result) {
					result["sqltext_digest"] = sql.Digest(result["sqltext"])
				}

				// Convert timestamp ISO 8601 (UTC) to RFC 3339:
				result["_time"] = cast.ToDateTime(result["timestamp"], "2006-01-02T15:04:05 UTC")
				result["host_ip"] = cnf.IPAddress
				result["host_name"] = cnf.General.Hostname
				result["sqltext"] = strings.Escape(result["sqltext"])
				result["sqltext_digest"] = strings.Escape(result["sqltext_digest"])

				// Remove unnused key:
				delete(result, "timestamp")
				delete(result, "record")

				parser <- result
			} else if line != "/>" {
				buffer = append(buffer, line)
			}
		}
	}()
}
