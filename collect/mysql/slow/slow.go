package slow

import (
  "regexp"
  "gitlab.com/swapbyt3s/zenit/common"
)

func Parser(path string, tail <-chan string, parser chan<- map[string]string) {
  var buffer string
  reRecord := regexp.MustCompile(`# Time: (.+);(.+);$`)
  reKeyVal := regexp.MustCompile(`\W+Time:\W+(?P<time>\d{6}\W\d{2}:\d{2}:\d{2})` +
                                 `\W+User@Host:\W+(?P<user_host>\w+\[\w+]\W+@\W+\[\d{1,3}.\d{1,3}.\d{1,3}.\d{1,3}\])`+
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
                                 `\W+SET timestamp=(?P<timestamp>\d+);`+
                                 `\W+(?P<query>.*);`)

  go func() {
    defer close(parser)

    for line := range tail {
      buffer += " " + line
      record := reRecord.FindString(buffer)

      if len(record) > 0 {
        buffer = ""
        parser <- common.RegexpGetGroups(reKeyVal, record)
      }
    }
  }()
}
