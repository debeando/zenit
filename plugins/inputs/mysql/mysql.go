package mysql

import (
  "log"
  "strings"

  "gitlab.com/swapbyt3s/zenit/common"
  "gitlab.com/swapbyt3s/zenit/config"
)

func Check() bool {
  log.Printf("I! - MySQL - DSN: %s\n", config.MySQL.DSN)
  conn, err := common.MySQLConnect(config.MySQL.DSN)
  if err != nil {
    log.Printf("E! - MySQL - Imposible to connect: %s\n", err)
    return false
  }

  log.Println("I! - MySQL - Connected successfully.")
  conn.Close()
  return true
}

func ClearUser(u string) string {
  index := strings.Index(u, "[")
  if index > 0 {
    return u[0:index]
  }
  return u
}
