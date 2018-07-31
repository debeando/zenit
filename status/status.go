package status

import (
  "fmt"

  "gitlab.com/swapbyt3s/zenit/config"
  "gitlab.com/swapbyt3s/zenit/common"
)

func Run() {
  CheckMySQL()
  CheckProxySQL()
  CheckClickHouse()
}

func CheckMySQL() {
  fmt.Printf("==> MySQL\n")
  fmt.Printf("--> Config: %s\n", config.DSN_MYSQL)
  conn, err := common.MySQLConnect(config.DSN_MYSQL)
  if err != nil {
    fmt.Printf("--> Status: Error\n")
    fmt.Printf("--> Message: %s\n", err)
  } else {
    fmt.Printf("--> Status: Ok\n")
    conn.Close()
  }
}

func CheckProxySQL() {
  fmt.Printf("==> ProxySQL\n")
  fmt.Printf("--> Config: %s\n", config.DSN_PROXYSQL)
  conn, err := common.MySQLConnect(config.DSN_PROXYSQL)
  if err != nil {
    fmt.Printf("--> Status: Error\n")
    fmt.Printf("--> Message: %s\n", err)
  } else {
    fmt.Printf("--> Status: Ok\n")
    conn.Close()
  }
}

func CheckClickHouse() {
  fmt.Printf("==> ClickHouse\n")
  fmt.Printf("--> Config: %s\n", config.DSN_CLICKHOUSE)
  if ! common.HTTPPost(config.DSN_CLICKHOUSE, "SELECT 1;") {
    fmt.Printf("--> Status: Error\n")
  } else {
    fmt.Printf("--> Status: Ok\n")
  }
}

func CheckProcess() {

}
