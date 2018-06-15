package mysql

import (
  "fmt"
)

func PrometheusExport() {
  tables := LoadTables()

  for i := range(tables.Items) {
    path := "mysql_stats_tables"
    path  = path + fmt.Sprintf("{schema=\"%s\"", tables.GetSchema(i))
    path  = path + fmt.Sprintf(",table=\"%s\"}",  tables.GetTable(i))
    fmt.Printf("%s %d\n", path, tables.GetSize(i))
  }
}
