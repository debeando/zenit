package mysql

import (
  "fmt"
  "strings"
  "gitlab.com/swapbyt3s/zenit/config"
  "gitlab.com/swapbyt3s/zenit/common"
)

type Column struct {
  schema    string
  table     string
  column    string
  data_type string
  unsigned  bool
  current   uint64
  percent   float64
}

type Columns struct {
  Items []Column
}

var list_columns *Columns

const (
  dt_tinyint      uint8  = 127
  dt_smallint     uint16 = 32767
  dt_mediumint    uint32 = 8388607
  dt_int          uint32 = 2147483647
  dt_bigint       uint64 = 9223372036854775807
  dt_us_tinyint   uint8  = 255
  dt_us_smallint  uint16 = 65535
  dt_us_mediumint uint32 = 16777215
  dt_us_int       uint32 = 4294967295
  dt_us_bigint    uint64 = 18446744073709551615
  QUERY_SQL_COLUMNS = `
SELECT table_schema, table_name, column_name, column_type
FROM information_schema.columns
WHERE table_schema NOT IN ('mysql','sys','performance_schema','information_schema','percona')
  AND column_type LIKE '%int%'
  AND column_key = 'PRI'
`
  QUERY_SQL_MAX_INT = "SELECT COALESCE(MAX(%s), 0) FROM %s.%s"
)

func LoadColumns() *Columns {
  if list_columns == nil {
    list_columns = &Columns{}
  }
  return list_columns
}

func (c *Columns) AddItem(item Column) []Column {
  c.Items = append(c.Items, item)
  return c.Items
}

func (c *Columns) GetSchema(i int) string {
  return c.Items[i].schema
}

func (c *Columns) GetTable(i int) string {
  return c.Items[i].table
}

func (c *Columns) GetPercent(i int) float64 {
  return c.Items[i].percent
}

func (c *Columns) GetUnsigned(i int) bool {
  return c.Items[i].unsigned
}

func (c *Columns) GetDataType(i int) string {
  return c.Items[i].data_type
}


func GatherOverflow() {
  conn, err := common.MySQLConnect(config.DSN_MYSQL)
  defer conn.Close()
  if err != nil {
    panic(err)
  }

  rows, err := conn.Query(QUERY_SQL_COLUMNS)
  defer rows.Close()
  if err != nil {
    panic(err)
  }

  columns := LoadColumns()

  for rows.Next() {
    var c Column
    var m uint64

    rows.Scan(
      &c.schema,
      &c.table,
      &c.column,
      &c.data_type)

    err = conn.QueryRow(fmt.Sprintf(QUERY_SQL_MAX_INT, c.column, c.schema, c.table)).Scan(&m)
    if err != nil {
      panic(err)
    }

    c.unsigned  = strings.Contains(c.data_type, "unsigned")
    c.data_type =c.data_type[0:strings.Index(c.data_type, "(")]
    c.current  = m

    if c.unsigned == true {
      switch c.data_type {
      case  "tinyint":
        c.percent  = (float64(c.current) / float64(dt_us_tinyint)) * 100
      case  "smallint":
        c.percent  = (float64(c.current) / float64(dt_us_smallint)) * 100
      case  "mediumint":
        c.percent  = (float64(c.current) / float64(dt_us_mediumint)) * 100
      case  "int":
        c.percent  = (float64(c.current) / float64(dt_us_int)) * 100
      case  "bigint":
        c.percent  = (float64(c.current) / float64(dt_us_bigint)) * 100
      }
    } else {
      switch c.data_type {
      case  "tinyint":
        c.percent  = (float64(c.current) / float64(dt_tinyint)) * 100
      case  "smallint":
        c.percent  = (float64(c.current) / float64(dt_smallint)) * 100
      case  "mediumint":
        c.percent  = (float64(c.current) / float64(dt_mediumint)) * 100
      case  "int":
        c.percent  = (float64(c.current) / float64(dt_int)) * 100
      case  "bigint":
        c.percent  = (float64(c.current) / float64(dt_bigint)) * 100
      }
    }

    columns.AddItem(c)
  }
}
