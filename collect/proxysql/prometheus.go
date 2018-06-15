package proxysql

import (
  "fmt"
)

func PrometheusExport() {
  queries := LoadQueries()
  queries.Sort()

  for i := range(queries.Items) {
    path := "proxysql_stats_queries"
    path  = path + fmt.Sprintf("{schema=\"%s\"",  queries.GetSchema(i))
    path  = path + fmt.Sprintf(",table=\"%s\"",   queries.GetTable(i))
    path  = path + fmt.Sprintf(",command=\"%s\"", queries.GetCommand(i))
    path  = path + fmt.Sprintf(",group=\"%s\"",   queries.GetGroup(i))

    fmt.Printf("%s,calc=\"count\"} %d\n", path, queries.GetCount(i))
    fmt.Printf("%s,calc=\"time\"} %d\n", path, queries.GetSum(i))
    fmt.Printf("%s,calc=\"min\"} %d\n", path, queries.GetMin(i))
    fmt.Printf("%s,calc=\"max\"} %d\n", path, queries.GetMax(i))
    fmt.Printf("%s,calc=\"avg\"} %d\n", path, queries.GetAvg(i))
  }
}
