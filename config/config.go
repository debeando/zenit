// Package to have a global variables.
package config

import (
  "os"
)

var VERSION          string = "0.1.0"
var DEBUG            bool   = false
var DSN_SRC_PROXYSQL string = "radminuser:radminpass@tcp(127.0.0.1:6032)/"
// var DSN_SRC_MYSQL    string = "root:@tcp(127.0.0.1:3306)/information_schema"
var DSN_DST_MYSQL    string = "root:@tcp(127.0.0.1:3306)/proxysql_stats"
var HOST_TOKEN       string

func init() {
  env_dsn_src_proxysql := os.Getenv("DSN_SRC_PROXYSQL")
  // env_dsn_src_mysql    := os.Getenv("DSN_SRC_MYSQL")
  env_dsn_dst_mysql    := os.Getenv("DSN_DST_MYSQL")

  if env_dsn_src_proxysql != "" {
    DSN_SRC_PROXYSQL = env_dsn_src_proxysql
  }

  if env_dsn_dst_mysql != "" {
    DSN_DST_MYSQL = env_dsn_dst_mysql
  }
}
