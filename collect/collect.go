package collect

import (
  "gitlab.com/swapbyt3s/zenit/common"
  "gitlab.com/swapbyt3s/zenit/output"
  "gitlab.com/swapbyt3s/zenit/collect/mysql"
  "gitlab.com/swapbyt3s/zenit/collect/os"
  "gitlab.com/swapbyt3s/zenit/collect/percona"
  "gitlab.com/swapbyt3s/zenit/collect/proxysql"
)

func Run(services []string) {
  // OS
  if common.StringInArray("os", services) {
    os.GatherSysLimits()
    os.GatherMem()
    os.GatherCPU()
  }
  if common.StringInArray("os-limits", services) {
    os.GatherSysLimits()
  }
  if common.StringInArray("os-mem", services) {
    os.GatherMem()
  }
  if common.StringInArray("os-cpu", services) {
    os.GatherCPU()
  }

  // MySQL
  if common.StringInArray("mysql", services) {
    mysql.GatherStatus()
    mysql.GatherTables()
    mysql.GatherOverflow()
  }
  if common.StringInArray("mysql-status", services) {
    mysql.GatherStatus()
  }
  if common.StringInArray("mysql-tables", services) {
    mysql.GatherTables()
  }
  if common.StringInArray("mysql-overflow", services) {
    mysql.GatherOverflow()
  }

  mysql.PrometheusExport()

  // Percona
  if common.StringInArray("percona-process", services) {
    percona.GatherRunningProcess()
  }

  // ProxySQL
  if common.StringInArray("proxysql", services) {
    proxysql.GatherQueries()
  }

  proxysql.PrometheusExport()

  output.Prometheus()
}
