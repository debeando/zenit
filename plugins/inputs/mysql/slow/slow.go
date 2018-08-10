// TODO: Move this package to inputs/parsers/mysqlslowlog

package slow

import (
  "regexp"

  "github.com/swapbyt3s/zenit/plugins/inputs/mysql"
  "github.com/swapbyt3s/zenit/common"
  "github.com/swapbyt3s/zenit/config"
)

const (
  ROW = `#\W(Time|User\@Host):.+SET.+(SELECT|INSERT|UPDATE|DELETE).+;`
  KV  = `(\W+Time:\W+(?P<time>\d{6}\W\d{2}:\d{2}:\d{2}))?`+
        `\W+User@Host:\W+(?P<user_host>\w+\[\w+\]\W+@.*\[(\d{1,3}.\d{1,3}.\d{1,3}.\d{1,3})?\])`+
        `\W+Thread_id:\W+(?P<thread_id>\d+)`+
        `\W+Schema:\W+(?P<schema>\w+)`+
        `\W+Last_errno:\W+(?P<last_errno>\d+)`+
        `\W+Killed:\W+(?P<killed>\d+)`+
        `\W+Query_time:\W+(?P<query_time>\d+\.\d+)`+
        `\W+Lock_time:\W+(?P<lock_time>\d+\.\d+)`+
        `\W+Rows_sent:\W+(?P<rows_sent>\d+)`+
        `\W+Rows_examined:\W+(?P<rows_examined>\d+)`+
        `\W+Rows_affected:\W+(?P<rows_affected>\d+)`+
        `\W+Rows_read:\W+(?P<rows_read>\d+)`+
        `\W+Bytes_sent:\W+(?P<bytes_sent>\d+)`+
        `.+`+
        `\W+SET timestamp=(?P<timestamp>\d+);`+
        `\W+(?P<query>.*);`
)

func Parser(path string, tail <-chan string, parser chan<- map[string]string) {
  var buffer string

  go func() {
    defer close(parser)
    reRow := regexp.MustCompile(ROW)
    reKV  := regexp.MustCompile(KV)

    for line := range tail {
      buffer += line + " "
      record := reRow.FindString(buffer)

      if len(record) > 0 {
        buffer = ""
        result := common.RegexpGetGroups(reKV, record)

        if common.KeyInMap("user_host", result) {
          result["user_host"] = mysql.ClearUser(result["user_host"])
        }

        if common.KeyInMap("query", result) {
          result["query_digest"] = common.NormalizeQuery(result["query"])
        }

        result["host_ip"]      = config.IpAddress
        result["host_name"]    = config.General.Hostname
        result["query"]        = common.Escape(result["query"])
        result["query_digest"] = common.Escape(result["query_digest"])

        parser <- result
      }
    }
  }()
}
