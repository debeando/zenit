package main

import (
  "database/sql"
  "fmt"
  "time"

  _ "github.com/go-sql-driver/mysql"
)

func main() {
  db, err := sql.Open("mysql", "root@/mysql")
  if err != nil {
      panic(err.Error())
  }
  defer db.Close()

  tinit := time.Now()
  _, err = db.Query("SELECT 1") // This opens a connection first
  diff := time.Since(tinit)
  if err != nil {
    fmt.Println("Query1 err:", err.Error())
  } else {
    fmt.Println("Query1 time:", diff.String())
  }
}
