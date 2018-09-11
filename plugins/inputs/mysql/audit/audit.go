// This parse OLD Style
// https://dev.mysql.com/doc/refman/5.5/en/audit-log-file-formats.html
// TODO: Move this package to inputs/parsers/mysqlauditlog

package audit

import (
	"github.com/swapbyt3s/zenit/common"
	"github.com/swapbyt3s/zenit/common/mysql"
	"github.com/swapbyt3s/zenit/common/sql"
	"github.com/swapbyt3s/zenit/config"
)

func Parser(path string, tail <-chan string, parser chan<- map[string]string) {
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
					key, value := common.SplitKeyAndValue(&kv)
					result[key] = common.Trim(&value)
				}

				buffer = buffer[:0]

				if common.KeyInMap("user", result) {
					result["user"] = mysql.ClearUser(result["user"])
				}

				if common.KeyInMap("sqltext", result) {
					result["sqltext_digest"] = sql.Digest(result["sqltext"])
				}

				// Convert timestamp ISO 8601 (UTC) to RFC 3339:
				result["_time"] = common.ToDateTime(result["timestamp"], "2006-01-02T15:04:05 UTC")
				result["host_ip"] = config.IPAddress
				result["host_name"] = config.File.General.Hostname
				result["sqltext"] = common.Escape(result["sqltext"])
				result["sqltext_digest"] = common.Escape(result["sqltext_digest"])

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
