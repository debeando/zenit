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
)

func Parser(path string, tail <-chan string, parser chan<- map[string]string) {
  var buffer string
  reRecord := regexp.MustCompile(`<AUDIT_RECORD(.*?)/>`)
  reKeyVal := regexp.MustCompile(`(\w+)=("[^"]*")`)

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

        if user, ok := result["user"]; ok {
          result["user"] = clearUser(user)
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

func clearUser(u string) string {
  index := strings.Index(u, "[")
  if index > 0 {
    return u[0:index]
  }
  return u
}
