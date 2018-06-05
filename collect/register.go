// Package collect all stats from ProxySQL & MySQL and ingest into MySQL.
package collect

import (
  "fmt"
  "os"
  "github.com/swapbyt3s/zenit/config"
  "github.com/swapbyt3s/zenit/lib"
)

func Prepare() {
  fmt.Printf("==> Prepare...\n")

  _, err := lib.Connect(config.DSN_DST_MYSQL)
  if err == nil {
    fmt.Printf("--> Ping to MySQL: Ok\n")
  } else {
    fmt.Printf("--> Ping to MySQL: %s.\n", err)
    fmt.Printf("--> Imposible to continue.\n")
    os.Exit(1)
  }

  _, err = lib.Connect(config.DSN_SRC_PROXYSQL)
  if err == nil {
    fmt.Printf("--> Ping to ProxySQL: Ok\n")
  } else {
    fmt.Printf("--> Ping to ProxySQL: %s.\n", err)
    fmt.Printf("--> Imposible to continue.\n")
    os.Exit(1)
  }

  err = Register("ProxySQL", config.DSN_SRC_PROXYSQL)
  if err == nil {
    fmt.Printf("--> Registered ProxySQL: Ok\n")
  } else {
    fmt.Printf("--> Registered ProxySQL: %s\n", err)
  }

  err = Register("MySQL", config.DSN_DST_MYSQL)
  if err == nil {
    fmt.Printf("--> Registered MySQL: Ok\n")
  } else {
    fmt.Printf("--> Registered MySQL: %s\n", err)
  }

  config.HOST_TOKEN = lib.MD5(config.DSN_SRC_PROXYSQL)
}

func Register(name string, dsn string) error {
  sqlStatement := "INSERT IGNORE INTO hosts (id, dsn, role) VALUES(?, ?, ?)"
  conn, _ := lib.Connect(config.DSN_DST_MYSQL)
  _, err := conn.Exec(sqlStatement, lib.MD5(dsn), dsn, name)
  if err != nil {
    return err
  }
  return nil
}
