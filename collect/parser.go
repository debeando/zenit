package collect

import (
  "fmt"
  "gitlab.com/swapbyt3s/zenit/collect/mysql/audit"
  "gitlab.com/swapbyt3s/zenit/collect/mysql/slow"
  "gitlab.com/swapbyt3s/zenit/common"
  "gitlab.com/swapbyt3s/zenit/output/clickhouse"
)

func Parser(parse string, path string) {
  switch parse {
  case "auditlog-xml":
    AuditLogXML(path)
  case "slowlog":
    SlowLog(path)
  default:
    fmt.Println("Valid parser options:")
    fmt.Println("- auditlog-xml")
    fmt.Println("- slowlog")
  }
}

func AuditLogXML(path string) {
  channel_tail   := make(chan string)
  channel_parser := make(chan map[string]string)
  channel_event  := make(chan map[string]string)

  go common.Tail(path, channel_tail)
  go audit.Parser(path, channel_tail, channel_parser)
  go clickhouse.SendMySQLAuditLog(channel_event)

  for event := range channel_parser {
    channel_event <- event
  }
}

func SlowLog(path string) {
  channel_tail   := make(chan string)
  channel_parser := make(chan map[string]string)
  channel_event  := make(chan map[string]string)

  go common.Tail(path, channel_tail)
  go slow.Parser(path, channel_tail, channel_parser)
  go clickhouse.SendMySQLSlowLog(channel_event)

  for event := range channel_parser {
    channel_event <- event
  }
}
