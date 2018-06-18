package collect

import (
  "gitlab.com/swapbyt3s/zenit/lib"
  "gitlab.com/swapbyt3s/zenit/collect/os"
  "gitlab.com/swapbyt3s/zenit/collect/mysql"
  "gitlab.com/swapbyt3s/zenit/collect/percona"
  "gitlab.com/swapbyt3s/zenit/collect/proxysql"
)

func Run(services []string) {
  // OS
  if lib.StringInArray("os", services) {
    os.GatherSysLimits()
    os.GatherMem()
    os.GatherCPU()
  }
  if lib.StringInArray("os-limits", services) {
    os.GatherSysLimits()
  }
  if lib.StringInArray("os-mem", services) {
    os.GatherMem()
  }
  if lib.StringInArray("os-cpu", services) {
    os.GatherCPU()
  }

  // MySQL
  if lib.StringInArray("mysql", services) {
    mysql.GatherStatus()
    mysql.GatherTables()
    mysql.GatherOverflow()
  }
  if lib.StringInArray("mysql-status", services) {
    mysql.GatherStatus()
  }
  if lib.StringInArray("mysql-tables", services) {
    mysql.GatherTables()
  }
  if lib.StringInArray("mysql-overflow", services) {
    mysql.GatherOverflow()
  }

  mysql.PrometheusExport()

  // Percona
  if lib.StringInArray("percona-process", services) {
    percona.GatherRunningProcess()
  }

  // ProxySQL
  if lib.StringInArray("proxysql", services) {
    proxysql.GatherQueries()
  }

  proxysql.PrometheusExport()
}
