package collect

import (
  "gitlab.com/swapbyt3s/zenit/collect/os"
  "gitlab.com/swapbyt3s/zenit/collect/mysql"
  "gitlab.com/swapbyt3s/zenit/collect/percona"
  "gitlab.com/swapbyt3s/zenit/collect/proxysql"
)

func OS() {
  os.Run()
}

func MySQL() {
  mysql.Run()
}

func Percona() {
  percona.Run()
}

func ProxySQL() {
  proxysql.Run()
}
