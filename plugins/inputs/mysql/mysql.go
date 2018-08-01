package mysql

import (
  "log"
  "strings"

  "gitlab.com/swapbyt3s/zenit/common"
  "gitlab.com/swapbyt3s/zenit/config"
)

func Check() {
  log.Printf("MySQL - DSN: %s\n", config.DSN_MYSQL)
  conn, err := common.MySQLConnect(config.DSN_MYSQL)
  if err != nil {
    log.Printf("MySQL - Imposible to connect: %s\n", err)
  } else {
    log.Println("MySQL - Connected successfully.")
    conn.Close()
  }
}

func ClearUser(u string) string {
  index := strings.Index(u, "[")
  if index > 0 {
    return u[0:index]
  }
  return u
}
