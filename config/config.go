package config

import (
  "gitlab.com/swapbyt3s/zenit/common"
)

var(
  AUTHOR         string = "Nicola Strappazzon C. <swapbyt3s@gmail.com>"
  DSN_CLICKHOUSE string = "http://127.0.0.1:8123/?database=zenit"
  DSN_MYSQL      string = "root@tcp(127.0.0.1:3306)/"
  DSN_PROXYSQL   string = "radminuser:radminpass@tcp(127.0.0.1:6032)/"
  HOSTNAME       string = ""
  IPADDRESS      string = ""
  LOG_FILE       string = "/var/log/zenit.log"
  SLACK_CHANNEL  string = "alerts"
  SLACK_TOKEN    string = ""
  VERSION        string = "0.1.4"
)

func init() {
  DSN_CLICKHOUSE = common.GetEnv("DSN_CLICKHOUSE", DSN_CLICKHOUSE)
  DSN_MYSQL      = common.GetEnv("DSN_MYSQL", DSN_MYSQL)
  DSN_PROXYSQL   = common.GetEnv("DSN_PROXYSQL", DSN_PROXYSQL)
  SLACK_CHANNEL  = common.GetEnv("SLACK_CHANNEL", SLACK_CHANNEL)
  SLACK_TOKEN    = common.GetEnv("SLACK_TOKEN", "")
  IPADDRESS      = common.IpAddress()
  HOSTNAME       = common.Hostname()
}
