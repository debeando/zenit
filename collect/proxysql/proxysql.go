// Package collect all stats from ProxySQL and ingest into MySQL.
package proxysql

import (
  "github.com/swapbyt3s/zenit/collect/proxysql/dbmain"
  "github.com/swapbyt3s/zenit/collect/proxysql/dbstats"
)

func Run() {
  dbmain.Servers()
//  dbstats.Command()
  dbstats.Connections()
  dbstats.Queries()
}
