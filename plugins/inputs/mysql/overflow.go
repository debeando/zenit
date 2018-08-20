package mysql

import (
  "fmt"
  "strings"
  "strconv"

  "github.com/swapbyt3s/zenit/common"
  "github.com/swapbyt3s/zenit/config"
  "github.com/swapbyt3s/zenit/plugins/accumulator"
)

type Column struct {
  schema    string
  table     string
  column    string
  dataType string
  unsigned  bool
  current   uint64
  percent   float64
}

const (
  dtTinyInt     uint8  = 127
  dtSmallInt    uint16 = 32767
  dtMediumInt   uint32 = 8388607
  dtInt         uint32 = 2147483647
  dtBigInt      uint64 = 9223372036854775807
  dtUSTinyInt   uint8  = 255
  dtUSSmallInt  uint16 = 65535
  dtUSMediumInt uint32 = 16777215
  dtUSInt       uint32 = 4294967295
  dtUSBigInt    uint64 = 18446744073709551615
  QuerySQLColumns = `
SELECT table_schema, table_name, column_name, column_type
FROM information_schema.columns
WHERE table_schema NOT IN ('mysql','sys','performance_schema','information_schema','percona')
  AND column_type LIKE '%int%'
  AND column_key = 'PRI'
`
  QuerySQLMaxInt = "SELECT COALESCE(MAX(%s), 0) FROM %s.%s"
)

func Overflow() {
  conn, err := common.MySQLConnect(config.MySQL.DSN)
  if err != nil {
    panic(err.Error())
  }
  defer conn.Close()

  rows, err := conn.Query(QuerySQLColumns)
  if err != nil {
    panic(err.Error())
  }
  defer rows.Close()

  var a = accumulator.Load()

  for rows.Next() {
    var c Column
    var m uint64

    err := rows.Scan(
      &c.schema,
      &c.table,
      &c.column,
      &c.dataType)
    if err != nil {
      panic(err.Error())
    }

    err = conn.QueryRow(fmt.Sprintf(QuerySQLMaxInt, c.column, c.schema, c.table)).Scan(&m)
    if err != nil {
      panic(err)
    }

    c.unsigned  = strings.Contains(c.dataType, "unsigned")
    c.dataType =c.dataType[0:strings.Index(c.dataType, "(")]
    c.current  = m

    if c.unsigned {
      switch c.dataType {
      case  "tinyint":
        c.percent  = (float64(c.current) / float64(dtUSTinyInt)) * 100
      case  "smallint":
        c.percent  = (float64(c.current) / float64(dtUSSmallInt)) * 100
      case  "mediumint":
        c.percent  = (float64(c.current) / float64(dtUSMediumInt)) * 100
      case  "int":
        c.percent  = (float64(c.current) / float64(dtUSInt)) * 100
      case  "bigint":
        c.percent  = (float64(c.current) / float64(dtUSBigInt)) * 100
      }
    } else {
      switch c.dataType {
      case  "tinyint":
        c.percent  = (float64(c.current) / float64(dtTinyInt)) * 100
      case  "smallint":
        c.percent  = (float64(c.current) / float64(dtSmallInt)) * 100
      case  "mediumint":
        c.percent  = (float64(c.current) / float64(dtMediumInt)) * 100
      case  "int":
        c.percent  = (float64(c.current) / float64(dtInt)) * 100
      case  "bigint":
        c.percent  = (float64(c.current) / float64(dtBigInt)) * 100
      }
    }

    a.AddItem(accumulator.Metric{
      Key: "mysql_stats_overflow",
      Tags: []accumulator.Tag{accumulator.Tag{"schema", c.schema},
                              accumulator.Tag{"table", c.table},
                              accumulator.Tag{"type", "overflow"},
                              accumulator.Tag{"data_type", c.dataType},
                              accumulator.Tag{"unsigned", strconv.FormatBool(c.unsigned)}},
      Values: c.percent,
    })
  }
}
