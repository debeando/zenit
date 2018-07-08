package mysql

import (
  "database/sql"
  "gitlab.com/swapbyt3s/zenit/common"
  "gitlab.com/swapbyt3s/zenit/config"
  "gitlab.com/swapbyt3s/zenit/accumulator"
)

const QUERY_SQL_STATUS = "SHOW GLOBAL STATUS"

func GatherStatus() {
  conn, err := common.MySQLConnect(config.DSN_MYSQL)
  defer conn.Close()
  if err != nil {
    panic(err)
  }

  rows, err := conn.Query(QUERY_SQL_STATUS)
  defer rows.Close()
  if err != nil {
    panic(err)
  }

  var a = accumulator.Load()
  var k string
  var v sql.RawBytes

  for rows.Next() {
    rows.Scan(&k, &v)
    if value, ok := common.MySQLParseValue(v); ok {
      a.AddItem(accumulator.Metric{
        Key: "mysql_status",
        Tags: []accumulator.Tag{accumulator.Tag{"name", k}},
        Values: value,
      })
    }
  }
}
