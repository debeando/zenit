// Package generic lib for MySQL.
package lib

import (
  "log"
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

func ExecCollections(conn *sql.DB, statements []string) {
  for _, stmt := range statements {
    _, err := conn.Exec(stmt)
    if err != nil {
      log.Print(err.Error())
    }
  }
}
