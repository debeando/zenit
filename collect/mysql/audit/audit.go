// This parse OLD Style
// https://dev.mysql.com/doc/refman/5.5/en/audit-log-file-formats.html
// Todo:
// - Tener un mecanismo usando el random para obtener un % de muestra.
// - Tener un filtro para ignorar o pasar solo lo que nos interesa.
// - Tener un check para el logrotate en algun lado: /etc/logrotate.d/mysql.conf

package audit

import (
//  "fmt"
  "strings"
  "gitlab.com/swapbyt3s/zenit/collect/mysql"
  "gitlab.com/swapbyt3s/zenit/common"
)

var buffer []string

func Parser(path string, tail <-chan string, parser chan<- map[string]string) {
  go func() {
    defer close(parser)

    for line := range tail {
      if line == "<AUDIT_RECORD" && len(buffer) > 0 {
        result := make(map[string]string)

        for _, kv := range buffer {
          key, value := getKeyAndValue(kv)
          result[key] = Trim(value)
        }

        buffer = buffer[:0]

        if common.KeyInMap("user", result) {
          result["user"] = mysql.ClearUser(result["user"])
        }

        if common.KeyInMap("sqltext", result) {
          result["sqltext_digest"] = common.NormalizeQuery(result["sqltext"])
        }

        // For debug:
        // fmt.Printf("--(map)> %#v\n", result)

        parser <- result
      } else if line != "/>"{
        buffer = append(buffer, line)
      }
    }
  }()
}

func getKeyAndValue(s string) (key string, value string) {
  kv := strings.SplitN(s, "=", 2)
  if len(kv) == 2 {
    return strings.TrimSpace(strings.ToLower(kv[0])), kv[1]
  }
  return "", ""
}

func Trim(value string) string {
  value = strings.TrimSpace(value)
  value = strings.TrimRight(value, "\"")
  value = strings.TrimLeft(value, "\"")
  return value
}
