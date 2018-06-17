package mysql

import (
  "fmt"
)

func PrometheusExport() {
  status  := LoadStatus()
  tables  := LoadTables()
  columns := LoadColumns()

  for i := range(status.Items) {
    path := "mysql_status"
    path  = path + fmt.Sprintf("{name=\"%s\"}", status.GetName(i))
    fmt.Printf("%s %d\n", path, status.GetValue(i))
  }

  for i := range(tables.Items) {
    path := "mysql_stats_tables"
    path  = path + fmt.Sprintf("{schema=\"%s\"", tables.GetSchema(i))
    path  = path + fmt.Sprintf(",table=\"%s\"",  tables.GetTable(i))
    path  = path + fmt.Sprintf(",type=\"size\"}")
    fmt.Printf("%s %d\n", path, tables.GetSize(i))
  }

  for i := range(columns.Items) {
    path := "mysql_stats_tables"
    path  = path + fmt.Sprintf("{schema=\"%s\"", columns.GetSchema(i))
    path  = path + fmt.Sprintf(",table=\"%s\"",  columns.GetTable(i))
    path  = path + fmt.Sprintf(",type=\"overflow\"")
    path  = path + fmt.Sprintf(",data_type=\"%s\"",  columns.GetDataType(i))
    path  = path + fmt.Sprintf(",unsigned=\"%t\"}",  columns.GetUnsigned(i))
    fmt.Printf("%s %.2f\n", path, columns.GetPercent(i))
  }
}
