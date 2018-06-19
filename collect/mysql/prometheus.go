package mysql

import (
  "fmt"
)

func PrometheusExport() {
  tables  := LoadTables()

  for i := range(tables.Items) {
    path := "mysql_stats_tables"
    path  = path + fmt.Sprintf("{schema=\"%s\"", tables.GetSchema(i))
    path  = path + fmt.Sprintf(",table=\"%s\"",  tables.GetTable(i))
    path  = path + fmt.Sprintf(",type=\"size\"}")
    fmt.Printf("%s %d\n", path, tables.GetSize(i))

    path  = "mysql_stats_tables"
    path  = path + fmt.Sprintf("{schema=\"%s\"", tables.GetSchema(i))
    path  = path + fmt.Sprintf(",table=\"%s\"",  tables.GetTable(i))
    path  = path + fmt.Sprintf(",type=\"rows\"}")
    fmt.Printf("%s %d\n", path, tables.GetRows(i))

    path  = "mysql_stats_tables"
    path  = path + fmt.Sprintf("{schema=\"%s\"", tables.GetSchema(i))
    path  = path + fmt.Sprintf(",table=\"%s\"",  tables.GetTable(i))
    path  = path + fmt.Sprintf(",type=\"increment\"}")
    fmt.Printf("%s %d\n", path, tables.GetIncrement(i))
  }
}
