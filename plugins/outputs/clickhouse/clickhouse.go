// TODO: Add check for every X time to force send to CH and purge the buffer.

package clickhouse

import (
  "log"
  "time"

  "github.com/swapbyt3s/zenit/common"
  "github.com/swapbyt3s/zenit/common/sql"
  "github.com/swapbyt3s/zenit/config"
)

type Event struct {
  Type       string
  Schema     string
  Table      string
  Size       int
  Timeout    int
  Wildcard   map[string]string
  Values   []map[string]string
}

func Check() bool {
  log.Printf("I! - ClickHouse - DSN: %s\n", config.ClickHouse.DSN)
  if ! common.HTTPPost(config.ClickHouse.DSN, "SELECT 1;") {
    log.Println("E! - ClickHouse - Imposible to connect.")
    return false
  }

  log.Println("I! - ClickHouse - Connected successfully.")
  return true
}

func Run(e Event, data <-chan map[string]string, t int) {
  timeout := make(chan bool)
  ticker  := time.NewTicker(time.Second * 5)

  go func() {
    for range ticker.C {
      timeout <- true
    }
  }()

  for {
    select {
    case <-timeout:
      if config.General.Debug {
        log.Printf("D! - ClickHouse - Event timeout: %s - %#v\n", e.Type, e.Values)
      }
      if len(e.Values) > 0 {
        sql      := sql.Insert(e.Schema, e.Table, e.Wildcard, e.Values)
        e.Values = []map[string]string{}

        if config.General.Debug {
          log.Printf("D! - ClickHouse - Event insert: %s - %s", e.Type, sql)
        }

        go common.HTTPPost(config.ClickHouse.DSN, sql)
      }
    case d := <- data:
      if config.General.Debug {
        log.Printf("D! - ClickHouse - Event capture: %s - %#v\n", e.Type, d)
      }
      e.Values = append(e.Values, d)
      if len(e.Values) == e.Size {
        sql      := sql.Insert(e.Schema, e.Table, e.Wildcard, e.Values)
        e.Values = []map[string]string{}

        if config.General.Debug {
          log.Printf("D! - ClickHouse - Event insert: %s - %s", e.Type, sql)
        }

        go common.HTTPPost(config.ClickHouse.DSN, sql)
      }
    }
  }
}
