package mysql

import (
  "strings"
  "database/sql"
  "gitlab.com/swapbyt3s/zenit/config"
  "gitlab.com/swapbyt3s/zenit/common"
)

type State struct {
  name  string
  value uint64
}

type Status struct {
  Items []State
}

var list_status *Status

const QUERY_SQL_STATUS = "SHOW GLOBAL STATUS"

func LoadStatus() *Status {
  if list_status == nil {
    list_status = &Status{}
  }
  return list_status
}

func (s *Status) AddItem(item State) []State {
  s.Items = append(s.Items, item)
  return s.Items
}

func (s *Status) GetName(i int) string {
  return s.Items[i].name
}

func (s *Status) GetValue(i int) uint64 {
  return s.Items[i].value
}

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

  var key string
  var val sql.RawBytes

  status := LoadStatus()

  for rows.Next() {
    rows.Scan(&key, &val)

    if value, ok := common.MySQLParseValue(val); ok {
      status.AddItem(State{name: strings.ToLower(key), value: value})
    }
  }
}
