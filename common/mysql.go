package common

import (
  "bytes"
  "strconv"
  "database/sql"

  _ "github.com/go-sql-driver/mysql"
)

func MySQLConnect(dsn string) (*sql.DB, error) {
  db, err := sql.Open("mysql", dsn)
  if err != nil {
    return nil, err
  }

  err = db.Ping()
  if err != nil {
    return nil, err
  }

  return db, err
}

func MySQLParseValue(value sql.RawBytes) (uint64, bool) {
  if bytes.EqualFold(value, []byte("YES")) || bytes.Compare(value, []byte("ON")) == 0 {
    return 1, true
  }

  if bytes.EqualFold(value, []byte("NO")) || bytes.Compare(value, []byte("OFF")) == 0 {
    return 0, true
  }

  if val, err := strconv.ParseUint(string(value), 10, 64); err == nil {
    return val, true
  }

  return 0, false
}
