package proxysql

import (
  "github.com/swapbyt3s/zenit/common/mysql"
  "github.com/swapbyt3s/zenit/config"
)

func Check() bool {
  var enable = 0

  if ( config.File.ProxySQL.Inputs.Commands ) { enable++ }
  if ( config.File.ProxySQL.Inputs.Pool     ) { enable++ }
  if ( config.File.ProxySQL.Inputs.Queries  ) { enable++ }

  if enable > 0 {
    return mysql.Check(config.File.ProxySQL.DSN, "proxysql")
  }

  return false
}
