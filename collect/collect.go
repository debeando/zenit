// Package collect all stats from ProxySQL & MySQL and ingest into MySQL.
package collect

import (
  "fmt"
  "github.com/swapbyt3s/proxysql_crawler/collect/proxysql"
)

func Run() {
  fmt.Printf("==> Collect...\n")
  proxysql.Run()
}
