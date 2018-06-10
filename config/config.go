package config

import (
  "os"
)

var VERSION      string = "0.1.0"
var DSN_PROXYSQL string = "radminuser:radminpass@tcp(127.0.0.1:6032)/"

func init() {
  env_dsn_proxysql := os.Getenv("DSN_PROXYSQL")

  if env_dsn_proxysql != "" {
    DSN_PROXYSQL = env_dsn_proxysql
  }
}
