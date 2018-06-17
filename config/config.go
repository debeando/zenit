package config

import (
  "os"
)

var(
  VERSION      string = "0.1.0"
  DSN_MYSQL    string = "root@tcp(127.0.0.1:3306)/"
  DSN_PROXYSQL string = "radminuser:radminpass@tcp(127.0.0.1:6032)/"
)

func init() {
  env_dsn_mysql := os.Getenv("DSN_MYSQL")

  if env_dsn_mysql != "" {
    DSN_MYSQL = env_dsn_mysql
  }

  env_dsn_proxysql := os.Getenv("DSN_PROXYSQL")

  if env_dsn_proxysql != "" {
    DSN_PROXYSQL = env_dsn_proxysql
  }
}
