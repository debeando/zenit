package proxysql

import (
  "fmt"
)

func Prometheus() {
  stats := LoadStats()
  items := stats.ToArray()

  for i := range(items) {
   fmt.Printf("proxysql.stats.queries.%s\n", items[i])
  }
}
