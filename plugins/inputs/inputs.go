package inputs

import (
  "gitlab.com/swapbyt3s/zenit/common"
  "gitlab.com/swapbyt3s/zenit/plugins/inputs/mysql"
  "gitlab.com/swapbyt3s/zenit/plugins/inputs/os"
  "gitlab.com/swapbyt3s/zenit/plugins/inputs/percona"
  "gitlab.com/swapbyt3s/zenit/plugins/inputs/proxysql"
  "gitlab.com/swapbyt3s/zenit/plugins/outputs"
)

func Run(services []string) {
  // OS
  if common.StringInArray("os", services) {
    os.GatherCPU()
    os.GatherDisk()
    os.GatherMem()
    os.GatherNet()
    os.GatherSysLimits()
  }
  if common.StringInArray("os-cpu", services) {
    os.GatherCPU()
  }
  if common.StringInArray("os-disk", services) {
    os.GatherDisk()
  }
  if common.StringInArray("os-mem", services) {
    os.GatherMem()
  }
  if common.StringInArray("os-net", services) {
    os.GatherNet()
  }
  if common.StringInArray("os-limits", services) {
    os.GatherSysLimits()
  }

  // MySQL
  if common.StringInArray("mysql", services) {
    mysql.GatherOverflow()
    mysql.GatherSlave()
    mysql.GatherStatus()
    mysql.GatherTables()
    mysql.GatherVariables()
  }
  if common.StringInArray("mysql-overflow", services) {
    mysql.GatherOverflow()
  }
  if common.StringInArray("mysql-slave", services) {
    mysql.GatherSlave()
  }
  if common.StringInArray("mysql-status", services) {
    mysql.GatherStatus()
  }
  if common.StringInArray("mysql-tables", services) {
    mysql.GatherTables()
  }
  if common.StringInArray("mysql-variables", services) {
    mysql.GatherVariables()
  }

  // Percona
  if common.StringInArray("percona-process", services) {
    percona.GatherRunningProcess()
  }

  // ProxySQL
  if common.StringInArray("proxysql", services) {
    proxysql.GatherQueries()
  }

  // Output
  output.Prometheus()
}
