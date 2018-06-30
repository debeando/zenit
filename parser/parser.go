package parser

import (
  "fmt"
  "gitlab.com/swapbyt3s/zenit/collect/mysql/audit"
  "gitlab.com/swapbyt3s/zenit/common"
  "gitlab.com/swapbyt3s/zenit/config"
)

func Run(parse string, path string) {
  if parse == "auditlog-xml" {
    for record := range audit.Parse(path) {
      sql := audit.ToSQL(record)
      fmt.Printf("%s\n", sql)
      common.HTTPPost(config.CLICKHOUSE_API, sql)
    }
  }
}
