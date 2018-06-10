package lib

import (
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
)

func Connect(dsn string) (*sql.DB, error) {
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
