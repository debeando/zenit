// This parse OLD Style
// https://dev.mysql.com/doc/refman/5.5/en/audit-log-file-formats.html
// Todo:
// - Tener un mecanismo usando el random para obtener un % de muestra.
// - Tener un filtro para ignorar o pasar solo lo que nos interesa.
// - Tener un check para el logrotate en algun lado: /etc/logrotate.d/mysql.conf

package audit

import (
  "strings"
  "regexp"
  "gitlab.com/swapbyt3s/zenit/collect/mysql"
  "gitlab.com/swapbyt3s/zenit/common"
)

var (
  reRecord = regexp.MustCompile(`<AUDIT_RECORD(.*?)/>`)
  reKeyVal = regexp.MustCompile(`(\w+)=("[^"]*")`)
)

func Parser(path string, tail <-chan string, parser chan<- map[string]string) {
  var buffer string

  go func() {
    defer close(parser)

    for line := range tail {
      buffer += line
      record := reRecord.FindString(buffer)

      if len(record) > 0 {
        buffer = ""
        result := make(map[string]string)
        match := reKeyVal.FindAllString(record, -1)
        for i := range match {
          key, value := getKeyAndValue(match[i])
          value = trim(value)
          result[key] = value
        }

        if common.KeyInMap("user", result) {
          result["user"] = mysql.ClearUser(result["user"])
        }

        if common.KeyInMap("sqltext", result) {
          result["sqltext_digest"] = common.NormalizeQuery(result["sqltext"])
        }

        parser <- result
      }
    }
  }()
}

func getKeyAndValue(s string) (key string, value string) {
  kv := strings.SplitN(s, "=", 2)
  return strings.ToLower(kv[0]), kv[1]
}

func trim(value string) string {
  value = strings.TrimRight(value, "\"")
  value = strings.TrimLeft(value, "\"")
  return value
}
