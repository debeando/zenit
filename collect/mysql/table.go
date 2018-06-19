package mysql

import (
  "gitlab.com/swapbyt3s/zenit/config"
  "gitlab.com/swapbyt3s/zenit/common"
)

type Table struct {
  schema    string
  table     string
  size      uint64
  rows      uint64
  increment uint64
}

type Tables struct {
  Items []Table
}

var list_tables *Tables

const QUERY_SQL_TABLES = `
SELECT table_schema AS 'schema',
       table_name AS 'table',
       data_length + index_length AS 'size',
       table_rows AS 'rows',
       auto_increment AS 'increment'
FROM information_schema.tables
WHERE table_schema NOT IN ('mysql','sys','performance_schema','information_schema','percona')
ORDER BY table_schema, table_name;
`

func LoadTables() *Tables {
  if list_tables == nil {
    list_tables = &Tables{}
  }
  return list_tables
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

func (tables *Tables) GetRows(i int) uint64 {
  return tables.Items[i].rows
}

func (tables *Tables) GetIncrement(i int) uint64 {
  return tables.Items[i].increment
}

func GatherTables() {
  conn, err := common.MySQLConnect(config.DSN_MYSQL)
  defer conn.Close()
  if err != nil {
    panic(err)
  }

  rows, err := conn.Query(QUERY_SQL_TABLES)
  defer rows.Close()
  if err != nil {
    panic(err)
  }

  tables := LoadTables()

  for rows.Next() {
    var t Table

    rows.Scan(
      &t.schema,
      &t.table,
      &t.size,
      &t.rows,
      &t.increment)

    tables.AddItem(t)
  }
}
