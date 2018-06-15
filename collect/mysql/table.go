package mysql

import (
  "gitlab.com/swapbyt3s/zenit/config"
  "gitlab.com/swapbyt3s/zenit/lib"
)

type Table struct {
  schema string
  table  string
  size   uint64
}

type Tables struct {
  Items []Table
}

var list *Tables

const QUERY_SQL = `
SELECT table_schema AS 'schema', table_name AS 'table', data_length + index_length AS 'size'
FROM information_schema.tables
WHERE table_schema NOT IN ('mysql', 'sys', 'performance_schema', 'information_schema')
ORDER BY table_schema, table_name;
`

func LoadTables() *Tables {
  if list == nil {
    list = &Tables{}
  }
  return list
}

func (tables *Tables) AddItem(item Table) []Table {
  tables.Items = append(tables.Items, item)
  return tables.Items
}

func (tables *Tables) GetSchema(i int) string {
  return tables.Items[i].schema
}

func (tables *Tables) GetTable(i int) string {
  return tables.Items[i].table
}

func (tables *Tables) GetSize(i int) uint64 {
  return tables.Items[i].size
}

func GetTables() {
  conn, err := lib.MySQLConnect(config.DSN_PROXYSQL)
  defer conn.Close()
  if err != nil {
    panic(err)
  }

  rows, err := conn.Query(QUERY_SQL)
  defer rows.Close()
  if err != nil {
    panic(err)
  }

  tables := LoadTables()

  for rows.Next() {
    var table Table

    rows.Scan(
      &table.schema,
      &table.table,
      &table.size)

    tables.AddItem(table)
  }
}
