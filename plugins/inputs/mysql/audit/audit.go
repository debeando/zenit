// This parse OLD Style
// https://dev.mysql.com/doc/refman/5.5/en/audit-log-file-formats.html

package audit

import (
//  "fmt"
  "strings"

  "gitlab.com/swapbyt3s/zenit/plugins/inputs/mysql"
  "gitlab.com/swapbyt3s/zenit/common"
  "gitlab.com/swapbyt3s/zenit/config"
)

func Parser(path string, tail <-chan string, parser chan<- map[string]string) {
  var buffer []string

  go func() {
    defer close(parser)

    for line := range tail {
      if line == "<AUDIT_RECORD" && len(buffer) > 0 {
        result := make(map[string]string)

        for _, kv := range buffer {
          key, value := ParseKeyAndValue(&kv)
          result[key] = Trim(&value)
        }

        buffer = buffer[:0]

        if common.KeyInMap("user", result) {
          result["user"] = mysql.ClearUser(result["user"])
        }

        if common.KeyInMap("sqltext", result) {
          result["sqltext_digest"] = common.NormalizeQuery(result["sqltext"])
        }

        // Convert timestamp ISO 8601 (UTC) to RFC 3339:
        result["timestamp"]      = common.ToDateTime(result["timestamp"], "2006-01-02T15:04:05 UTC")
        result["host_ip"]        = config.IPADDRESS
        result["host_name"]      = config.HOSTNAME
        result["sqltext"]        = common.Escape(result["sqltext"])
        result["sqltext_digest"] = common.Escape(result["sqltext_digest"])

        // For debug:
        // fmt.Printf("--(map)> %#v\n", result)

        parser <- result
      } else if line != "/>"{
        buffer = append(buffer, line)
      }
    }
  }()
}

func ParseKeyAndValue(s *string) (key string, value string) {
  kv := strings.SplitN(*s, "=", 2)
  if len(kv) == 2 {
    return strings.TrimSpace(strings.ToLower(kv[0])), kv[1]
  }
  return "", ""
}

func Trim(value *string) string {
  *value = strings.TrimSpace(*value)
  *value = strings.TrimRight(*value, "\"")
  *value = strings.TrimLeft(*value, "\"")
  return *value
}
