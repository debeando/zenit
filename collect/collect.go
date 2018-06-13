package collect

import (
  "gitlab.com/swapbyt3s/zenit/collect/os"
  "gitlab.com/swapbyt3s/zenit/collect/proxysql"
)

func ProxySQL() {
  proxysql.Run()
}

func OS() {
  os.Run()
}
