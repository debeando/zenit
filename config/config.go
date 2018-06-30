package config

import (
  "gitlab.com/swapbyt3s/zenit/common"
)

var(
  CLICKHOUSE_API string = ""
  DSN_MYSQL      string = "root@tcp(127.0.0.1:3306)/"
  DSN_PROXYSQL   string = "radminuser:radminpass@tcp(127.0.0.1:6032)/"
  SLACK_CHANNEL  string = "alerts"
  SLACK_TOKEN    string = ""
  VERSION        string = "0.1.4"
)

func init() {
  CLICKHOUSE_API = common.GetEnv("CLICKHOUSE_API", CLICKHOUSE_API)
  SLACK_CHANNEL  = common.GetEnv("SLACK_CHANNEL", SLACK_CHANNEL)
  SLACK_TOKEN    = common.GetEnv("SLACK_TOKEN", "")
  DSN_MYSQL      = common.GetEnv("DSN_MYSQL", DSN_MYSQL)
  DSN_PROXYSQL   = common.GetEnv("DSN_PROXYSQL", DSN_PROXYSQL)
}
