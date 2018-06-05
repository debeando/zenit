// Package collect all stats from ProxySQL and ingest into MySQL.
package proxysql

import (
  "github.com/swapbyt3s/proxysql_crawler/collect/proxysql/dbmain"
  "github.com/swapbyt3s/proxysql_crawler/collect/proxysql/dbstats"
)

func Run() {
  dbmain.Servers()
//  dbstats.Command()
  dbstats.Connections()
  dbstats.Queries()
}
